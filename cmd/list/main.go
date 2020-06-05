package main

import (
	"flag"
	"fmt"
	"knative.dev/eventing/test/pkg/collector"
	"knative.dev/eventing/test/pkg/injection"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "knative.dev/pkg/system/testing"
)

func main() {
	// only scrape k8s.

	flag.Parse()
	ctx := injection.InjectionEnabled()

	c := collector.New(ctx)

	from := duckv1.KReference{
		Kind:       "Namespace",
		Name:       "default",
		APIVersion: "v1",
	}

	events, err := c.List(from)
	if err != nil {
		panic(err)
	}

	for i, e := range events {
		fmt.Printf("[%d]: seen by %q\n%s\n", i, e.Observer, e.Event)
	}
}
