package recorder

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	"knative.dev/eventing/test/pkg/api"
)

type recorder struct {
	out record.EventRecorder
	on  runtime.Object
}

func (r *recorder) Observe(observed api.Observed) error {
	b, err := json.Marshal(observed)
	if err != nil {
		return err
	}

	r.out.Eventf(r.on, corev1.EventTypeNormal, "EventObserved",
		"%s", string(b))

	return nil
}
