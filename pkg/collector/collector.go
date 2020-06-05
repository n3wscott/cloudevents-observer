package collector

import (
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"knative.dev/eventing/test/pkg/api"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
)

// Filter observed events, return true if the observation should be included.
type FilterFn func(ob api.Observed) bool

type Collector interface {
	// List all observed events from a ref, optionally filter (pass filter, && together)
	List(from duckv1.KReference, filters ...FilterFn) ([]api.Observed, error)
}

func New(ctx context.Context) Collector {
	return &collector{client: kubeclient.Get(ctx)}
}

type collector struct {
	client kubernetes.Interface
}

func (c *collector) List(from duckv1.KReference, filters ...FilterFn) ([]api.Observed, error) {
	var lister v1.EventInterface
	if from.Kind == "Namespace" {
		lister = c.client.CoreV1().Events(from.Name)
	} else {
		lister = c.client.CoreV1().Events(from.Namespace) // TODO: I do not understand how to do cluster scoped objects.
	}
	events, err := lister.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	obs := make([]api.Observed, 0)

	for _, v := range events.Items {
		switch v.Reason {
		case api.EventReason:
			ob := api.Observed{}
			if err := json.Unmarshal([]byte(v.Message), &ob); err != nil {
				return nil, err
			}
			if filters != nil {
				skip := false
				for _, fn := range filters {
					if !fn(ob) {
						skip = true
					}
				}
				if skip {
					continue
				}
			}

			obs = append(obs, ob)
			// TODO: worry about v.Count
		}
	}

	// sort by time collected
	return obs, nil
}
