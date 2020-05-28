package recorder

import (
	"context"
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"log"
	"os"

	"knative.dev/eventing/test/pkg/eventlog"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/logging"
)

type envConfig struct {
	AgentName string `envconfig:"AGENT_NAME" default:"observer-default" required:"true"`
	EventOn   string `envconfig:"K8S_EVENT_SINK" required:"true"`

	Port int    `envconfig:"PORT" default:"8080" required:"true"`
	Sink string `envconfig:"K_SINK"`
}

func NewFromEnv(ctx context.Context) eventlog.EventLog {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	var ref duckv1.KReference
	if err := json.Unmarshal([]byte(env.EventOn), &ref); err != nil {
		log.Printf("[ERROR] Failed to process env var [K8S_EVENT_SINK]: %s", err)
		os.Exit(1)
	}

	return NewEventLog(ctx, env.AgentName, ref)
}

func NewEventLog(ctx context.Context, agentName string, ref duckv1.KReference) eventlog.EventLog {

	gv, err := schema.ParseGroupVersion(ref.APIVersion)
	if err != nil {
		logging.FromContext(ctx).Fatalf("failed to parse group version, %s", err)
	}

	gvr, _ := meta.UnsafeGuessKindToResource(gv.WithKind(ref.Kind))

	var on runtime.Object
	if ref.Namespace == "" {
		on, err = dynamicclient.Get(ctx).Resource(gvr).Get(ref.Name, metav1.GetOptions{})
	} else {
		on, err = dynamicclient.Get(ctx).Resource(gvr).Get(ref.Name, metav1.GetOptions{})
	}
	if err != nil {
		logging.FromContext(ctx).Fatalf("failed to fetch object ref, %+v, %s", ref, err)

	}

	return &recorder{out: createRecorder(ctx, agentName), on: on}
}

func createRecorder(ctx context.Context, agentName string) record.EventRecorder {
	logger := logging.FromContext(ctx)

	recorder := controller.GetEventRecorder(ctx)
	if recorder == nil {
		// Create event broadcaster
		logger.Debug("Creating event broadcaster")
		eventBroadcaster := record.NewBroadcaster()
		watches := []watch.Interface{
			eventBroadcaster.StartLogging(logger.Named("event-broadcaster").Infof),
			eventBroadcaster.StartRecordingToSink(
				&v1.EventSinkImpl{Interface: kubeclient.Get(ctx).CoreV1().Events("")}),
		}
		recorder = eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: agentName})
		go func() {
			<-ctx.Done()
			for _, w := range watches {
				w.Stop()
			}
		}()
	}

	return recorder
}
