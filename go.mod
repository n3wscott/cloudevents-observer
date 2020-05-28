module knative.dev/eventing/test

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.0.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.17.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/utils v0.0.0-20200520001619-278ece378a50 // indirect
	knative.dev/pkg v0.0.0-20200528172301-6eb9fa1546b8
)

replace (
	k8s.io/api => k8s.io/api v0.17.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/code-generator => k8s.io/code-generator v0.17.6
)
