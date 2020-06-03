package cloudout

import (
	"context"
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/kelseyhightower/envconfig"
	"knative.dev/eventing/test/pkg/api"
	"knative.dev/eventing/test/pkg/eventlog"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"log"
	"os"
)

type envConfig struct {
	Sink string `envconfig:"K_SINK"`
	// CEOverrides are the CloudEvents overrides to be applied to the outbound event.
	CEOverrides string `envconfig:"K_CE_OVERRIDES"`
}

func NewFromEnv(ctx context.Context) eventlog.EventLog {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	if env.Sink == "" {
		return nil
	}

	p, err := cloudevents.NewHTTP(http.WithTarget(env.Sink))
	if err != nil {
		log.Printf("[ERROR] Failed to create cloudevents http protocol: %s", err)
		os.Exit(1)
	}

	ce, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Printf("[ERROR] Failed to create cloudevents client: %s", err)
		os.Exit(1)
	}

	var ceOverrides *duckv1.CloudEventOverrides
	if len(env.CEOverrides) > 0 {
		overrides := duckv1.CloudEventOverrides{}
		err := json.Unmarshal([]byte(env.CEOverrides), &overrides)
		if err != nil {
			log.Printf("[ERROR] Unparseable CloudEvents overrides %s: %v", env.CEOverrides, err)
			os.Exit(1)
		}
		ceOverrides = &overrides
	}

	return &cloudevent{out: ce, ceOverrides: ceOverrides}
}

type cloudevent struct {
	out         cloudevents.Client
	ceOverrides *duckv1.CloudEventOverrides
}

// Forwards the event.
func (w *cloudevent) Observe(observed api.Observed) error {
	event := observed.Event

	if w.ceOverrides != nil && w.ceOverrides.Extensions != nil {
		for n, v := range w.ceOverrides.Extensions {
			event.SetExtension(n, v)
		}
	}

	if result := w.out.Send(context.Background(), event); !cloudevents.IsACK(result) {
		return result
	}

	return nil
}

// this would be a observed event.
//func (w *cloudevent) Observe(observed api.Observed) error {
//	event := cloudevents.NewEvent()
//	event.SetSource(observed.Observer)
//	event.SetType("knative-eventing.observed")
//
//	if w.ceOverrides != nil && w.ceOverrides.Extensions != nil {
//		for n, v := range w.ceOverrides.Extensions {
//			event.SetExtension(n, v)
//		}
//	}
//
//	if err := event.SetData(cloudevents.ApplicationJSON, observed); err != nil {
//		return err
//	}
//
//	if result := w.out.Send(context.Background(), event); !cloudevents.IsACK(result) {
//		return result
//	}
//
//	return nil
//}
