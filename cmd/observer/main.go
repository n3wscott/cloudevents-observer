package main

import (
	"context"
	"fmt"
	"knative.dev/eventing/test/pkg/eventlog/recorder"
	"os"

	"knative.dev/eventing/test/pkg/eventlog"
	"knative.dev/eventing/test/pkg/eventlog/writer"
	"knative.dev/eventing/test/pkg/observer"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
)

func main() {
	ctx := signals.NewContext()
	cfg := sharedmain.ParseAndGetConfigOrDie()
	ctx, informers := injection.Default.SetupInformers(ctx, cfg)

	// Start the injection clients and informers.
	go func(ctx context.Context) {
		if err := controller.StartInformers(ctx.Done(), informers...); err != nil {
			panic(fmt.Sprintf("Failed to start informers - %s", err))
		}
		<-ctx.Done()
	}(ctx)

	logs := writer.NewEventLog(ctx, os.Stdout)
	events := recorder.NewFromEnv(ctx)

	obs := observer.Observer{EventLogs: []eventlog.EventLog{logs, events}, ID: "demo"}

	if err := obs.Start(ctx); err != nil {
		panic(err)
	}
}
