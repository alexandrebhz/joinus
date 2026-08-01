package monitoring

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewApplication creates a New Relic application from .env / process env
// (NEW_RELIC_ENABLED, NEW_RELIC_APP_NAME, NEW_RELIC_LICENSE_KEY, …).
// Returns nil when disabled or missing a license key.
func NewApplication() (*newrelic.Application, error) {
	enabled, _ := strconv.ParseBool(os.Getenv("NEW_RELIC_ENABLED"))
	if !enabled || os.Getenv("NEW_RELIC_LICENSE_KEY") == "" {
		return nil, nil
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(c *newrelic.Config) {
			c.ErrorCollector.IgnoreStatusCodes = []int{400, 401, 403, 404, 405, 422}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("new relic: %w", err)
	}

	_ = app.WaitForConnection(5 * time.Second)
	return app, nil
}

// Shutdown flushes buffered data. No-op when app is nil.
func Shutdown(app *newrelic.Application) {
	if app == nil {
		return
	}
	app.Shutdown(10 * time.Second)
}
