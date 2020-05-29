package observer

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"knative.dev/eventing/test/pkg/api"
	"knative.dev/eventing/test/pkg/eventlog"
	"time"
)

type Observer struct {
	Observer  string
	EventLogs []eventlog.EventLog
}

func New(eventLogs ...eventlog.EventLog) *Observer {
	return &Observer{
		EventLogs: eventLogs,
	}
}

type envConfig struct {
	Observer string `envconfig:"OBSERVER" default:"observer-default" required:"true"`
}

func (o *Observer) Start(ctx context.Context) error {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		return err
	}
	o.Observer = env.Observer

	ce, err := cloudevents.NewDefaultClient()
	if err != nil {
		return err
	}

	return ce.StartReceiver(ctx, o.OnEvent)
}

func (o *Observer) OnEvent(event cloudevents.Event) {
	obs := api.Observed{
		Event:    event,
		Origin:   "http://origin", // TODO: we do not have this part at the moment.
		Observer: o.Observer,
		Time:     cloudevents.Timestamp{time.Now()}.String(),
	}

	for _, el := range o.EventLogs {
		if el == nil {
			continue
		}
		if err := el.Observe(obs); err != nil {
			fmt.Println(err)
		}
	}
}
