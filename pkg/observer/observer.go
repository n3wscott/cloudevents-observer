package observer

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"knative.dev/eventing/test/pkg/api"
	"knative.dev/eventing/test/pkg/eventlog"
	"time"
)

type Observer struct {
	ID        string
	EventLogs []eventlog.EventLog
}

func (o *Observer) Start(ctx context.Context) error {
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
		Observer: o.ID,
		Time:     cloudevents.Timestamp{time.Now()}.String(),
	}

	for _, el := range o.EventLogs {
		if err := el.Observe(obs); err != nil {
			fmt.Println(err)
		}
	}
}
