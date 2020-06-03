package api

import cloudevents "github.com/cloudevents/sdk-go/v2"

const (
	EventReason = "EventObserved"
)

type Observed struct {
	Event    cloudevents.Event `json:"event"`
	Origin   string            `json:"origin"`
	Observer string            `json:"observer"`
	Time     string            `json:"time"`
}
