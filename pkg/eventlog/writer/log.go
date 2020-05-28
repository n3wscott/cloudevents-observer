package writer

import (
	"context"
	"encoding/json"
	"io"

	"knative.dev/eventing/test/pkg/api"
	"knative.dev/eventing/test/pkg/eventlog"
)

func NewEventLog(ctx context.Context, out io.Writer) eventlog.EventLog {
	return &writer{out: out}
}

type writer struct {
	out io.Writer
}

var newline = []byte("\n")

func (w *writer) Observe(observed api.Observed) error {
	b, err := json.Marshal(observed)
	if err != nil {
		return err
	}
	if _, err := w.out.Write(b); err != nil {
		return err
	}
	if _, err := w.out.Write(newline); err != nil {
		return err
	}

	return nil
}
