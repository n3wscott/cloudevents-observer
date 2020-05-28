package eventlog

import "knative.dev/eventing/test/pkg/api"

type EventLog interface {
	Observe(observed api.Observed) error
}
