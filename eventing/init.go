package eventing

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	eventingclient "knative.dev/eventing/pkg/client/clientset/versioned"
)

var knClientset *eventingclient.Clientset

func init() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	knClientset, err = eventingclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
