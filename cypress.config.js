const { defineConfig } = require('cypress')

module.exports = defineConfig({
  chromeWebSecurity: false,
  e2e: {
    setupNodeEvents(on, config) {
      on("before:browser:launch", (browser = {}, launchOptions) => {
        if (browser.name === "chrome") {
          launchOptions.args.push("--disable-dev-shm-usage");
        }
        return launchOptions;
      });
    },
    baseUrl: 'http://localhost:8888/',
    supportFile: false,
  },
})
