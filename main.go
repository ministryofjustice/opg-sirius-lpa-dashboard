package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/server"
	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()
	logger := telemetry.NewLogger("opg-sirius-lpa-dashboard")

	if err := run(ctx, logger); err != nil {
		logger.Error("fatal startup error", slog.Any("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	port := env.Get("PORT", "8080")
	webDir := env.Get("WEB_DIR", "web")
	siriusURL := env.Get("SIRIUS_URL", "http://localhost:9001")
	siriusPublicURL := env.Get("SIRIUS_PUBLIC_URL", "")
	prefix := env.Get("PREFIX", "")
	exportTraces := env.Get("TRACING_ENABLED", "0") == "1"

	layouts, _ := template.
		New("").
		Funcs(map[string]interface{}{
			"join": func(sep string, items []string) string {
				return strings.Join(items, sep)
			},
			"contains": func(xs interface{}, needle interface{}) bool {
				switch need := needle.(type) {
				case string:
					for _, x := range xs.([]string) {
						if x == need {
							return true
						}
					}

				case int:
					for _, x := range xs.([]int) {
						if x == need {
							return true
						}
					}
				}

				return false
			},
			"prefix": func(s string) string {
				return prefix + s
			},
			"sirius": func(s string) string {
				return siriusPublicURL + s
			},
			"formatDate": func(d interface{}) string {
				switch t := d.(type) {
				case time.Time:
					return t.Format("02 Jan 2006")
				case sirius.SiriusDate:
					return t.Format("02 Jan 2006")
				default:
					panic("can't format date")
				}
			},
			"isoDate": func(d time.Time) string {
				if d.IsZero() {
					return ""
				}

				return d.Format("2006-01-02")
			},
			"upper": func(s string) string {
				return strings.ToUpper(s)
			},
			"statusColour": func(s string) string {
				switch s {
				case "Perfect":
					return "green"
				case "Imperfect":
					return "red"
				case "Pending":
					return "blue"
				default:
					return "grey"
				}
			},
		}).
		ParseGlob(webDir + "/template/layout/*.gotmpl")

	files, _ := filepath.Glob(webDir + "/template/*.gotmpl")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}

	shutdown, err := telemetry.StartTracerProvider(ctx, logger, exportTraces)
	defer shutdown()
	if err != nil {
		return err
	}

	httpClient := http.DefaultClient
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)

	client, err := sirius.NewClient(httpClient, siriusURL)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           server.New(logger, client, tmpls, prefix, siriusURL, siriusPublicURL, webDir),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("listen and serve error", slog.Any("err", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return server.Shutdown(tc)
}
