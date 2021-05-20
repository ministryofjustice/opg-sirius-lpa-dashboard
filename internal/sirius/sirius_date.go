package sirius

import (
	"strings"
	"time"
)

type SiriusDate struct {
	time.Time
}

func (sd *SiriusDate) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	strInput = strings.ReplaceAll(strInput, "\\/", "/")

	newTime, err := time.Parse("02/01/2006", strInput)
	if err != nil {
		return err
	}

	sd.Time = newTime
	return nil
}
