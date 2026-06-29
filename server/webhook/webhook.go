package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	logrus "github.com/teslamotors/fleet-telemetry/logger"
	"github.com/teslamotors/fleet-telemetry/telemetry"
)

// Dispatcher posts telemetry records to a GarageWatts Edge Function.
// It implements the telemetry.Producer interface.
type Dispatcher struct {
	url    string
	secret string
	client *http.Client
	logger *logrus.Logger
}

// NewDispatcher builds a Dispatcher from FLEET_TELEMETRY_WEBHOOK_* env vars.
func NewDispatcher(logger *logrus.Logger) *Dispatcher {
	return &Dispatcher{
		url:    os.Getenv("FLEET_TELEMETRY_WEBHOOK_URL"),
		secret: os.Getenv("FLEET_TELEMETRY_WEBHOOK_SECRET"),
		client: &http.Client{Timeout: 10 * time.Second},
		logger: logger,
	}
}

// Produce posts a single telemetry record to the configured webhook.
func (d *Dispatcher) Produce(record *telemetry.Record) {
	payload, err := json.Marshal(record)
	if err != nil {
		d.logger.ErrorLog("webhook_marshal_error", err, nil)
		return
	}

	req, err := http.NewRequest("POST", d.url, bytes.NewBuffer(payload))
	if err != nil {
		d.logger.ErrorLog("webhook_request_error", err, nil)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.secret))

	resp, err := d.client.Do(req)
	if err != nil {
		d.logger.ErrorLog("webhook_dispatch_error", err, nil)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		d.logger.ErrorLog("webhook_unexpected_status", nil, logrus.LogInfo{"status": resp.StatusCode})
	}
}

// Close satisfies telemetry.Producer; the dispatcher holds no resources to release.
func (d *Dispatcher) Close() error { return nil }

// ProcessReliableAck satisfies telemetry.Producer; webhooks are fire-and-forget.
func (d *Dispatcher) ProcessReliableAck(entry *telemetry.Record) {}

// ReportError satisfies telemetry.Producer.
func (d *Dispatcher) ReportError(message string, err error, logInfo logrus.LogInfo) {
	d.logger.ErrorLog(message, err, logInfo)
}
