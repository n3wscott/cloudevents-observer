package bdd

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"knative.dev/eventing/test/pkg/api"
	"knative.dev/eventing/test/pkg/collector"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages-go/v10"
)

var currentEvent *cloudevents.Event

func CloudEventsFeatureContext(ctx context.Context, s *godog.Suite) {

	runID := uuid.New().String()

	s.Step(`^v1.0 CloudEvents Attributes:$`, func(attributes *messages.PickleStepArgument_PickleTable) error {
		event := cloudevents.NewEvent("1.0")

		event.SetExtension("runid", runID)

		for _, row := range attributes.Rows {
			key := row.Cells[0].Value
			value := row.Cells[1].Value

			switch key {
			case "key":
				// ignore the header
				continue
			case "id":
				event.SetID(value)
			case "type":
				event.SetType(value)
			case "source":
				event.SetSource(value)
			case "subject":
				event.SetSubject(value)
			case "dataschema":
				event.SetDataSchema(value)
			case "time":
				t, err := time.Parse(time.RFC3339, value)
				if err != nil {
					return err
				}
				event.SetTime(t)
			default:
				event.SetExtension(key, value)
			}
		}

		if err := event.Validate(); err != nil {
			return err
		}

		currentEvent = &event

		return nil
	})

	s.Step(`^JSON Data:$`, func(jsonData *messages.PickleStepArgument_PickleDocString) error {
		currentEvent.SetDataContentType(cloudevents.ApplicationJSON)
		currentEvent.DataEncoded = []byte(jsonData.Content)
		return nil
	})

	s.Step(`^the consumer is ready$`, func() error {
		// TODO: check this.

		return nil
	})

	s.Step(`^the event is sent to "([^"]*)"$`, func(consumer string) error {
		ctx := cloudevents.ContextWithTarget(context.Background(), consumer)

		client, err := cloudevents.NewDefaultClient()
		if err != nil {
			return err
		}

		if result := client.Send(ctx, *currentEvent); !cloudevents.IsACK(result) {
			return result
		}

		return nil
	})

	s.Step(`^the consumer got the event$`, func() error {

		time.Sleep(4 * time.Second) // Coldstart...

		c := collector.New(ctx)

		from := duckv1.KReference{
			Kind:       "Namespace",
			Name:       "default",
			APIVersion: "v1",
		}

		events, err := c.List(from, func(ob api.Observed) bool {
			gotID, err := types.ToString(ob.Event.Extensions()["runid"])
			if err != nil {
				return false
			}
			return gotID == runID
		})
		if err != nil {
			panic(err)
		}

		if len(events) != 1 {
			return fmt.Errorf("fail: expected exactly 1 event, got %d", len(events))
		}

		//for i, e := range events {
		//	fmt.Printf("[%d]: seen by %q\n%s\n", i, e.Observer, e.Event)
		//}

		got := events[0].Event.String()
		want := currentEvent.String()

		if diff := cmp.Diff(got, want); diff != "" {
			return fmt.Errorf("unexpected event (-want, +got) = %v", diff)
		}

		return nil
	})
}
