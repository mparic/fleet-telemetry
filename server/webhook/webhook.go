package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	logrus "github.com/sirupsen/logrus"
	"github.com/teslamotors/fleet-telemetry/datastore/simple"
	"github.com/teslamotors/fleet-telemetry/messages"
	"github.com/teslamotors/fleet-telemetry/server/airbrake"
	"github.com/teslamotors/fleet-telemetry/telemetry"
)

// Dispatcher posts telemetry events to a GarageWatts Edge Function
type Dispatcher struct {
	url    string
	secret string
	client *http.Client
	logger *logrus.Logger
}

func NewDispatcher(logger *logrus.Logger) *Dispatcher {
	return &Dispatcher{
		url:    os.Getenv("FLEET_TELEMETRY_WEBHOOK_URL"),
		secret: os.Getenv("FLEET_TELEMETRY_WEBHOOK_SECRET"),
		client: &http.Client{Timeout: 10 * time.Second},
		logger: logger,
	}
}

func (d *Dispatcher) Dispatch(record *telemetry.Record) {
	payload, err := json.Marshal(record)
	if err != nil {
		d.logger.Errorf("[webhook] marshal error: %v", err)
		return
	}

	req, err := http.NewRequest("POST", d.url, bytes.NewBuffer(payload))
	if err != nil {
		d.logger.Errorf("[webhook] request error: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.secret))

	resp, err := d.client.Do(req)
	if err != nil {
		d.logger.Errorf("[webhook] dispatch error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		d.logger.Errorf("[webhook] unexpected status: %d", resp.StatusCode)
	}
}

func (d *Dispatcher) ReportError(err error, r *telemetry.Record, logInfo simple.LogInfo) {}

func (d *Dispatcher) ProcessingStats() airbrake.Stats { return airbrake.Stats{} }
