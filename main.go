package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/qhh0205/kncrud/serving"
	"k8s.io/client-go/tools/clientcmd"
	knclientv1 "knative.dev/client/pkg/serving/v1"
	servingv1client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

var client *servingv1client.ServingV1Client
var knclient knclientv1.KnServingClient

func init() {
	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "/Users/hello/.kube/config", "path to Kubernetes config file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, _ = servingv1client.NewForConfig(config)
	knclient = knclientv1.NewKnServingClient(client, "default")
}

func main() {
	ksvcCfg := serving.ServiceConfiguration{
		Name:              "jim",
		Namespace:         "default",
		Image:             "nginx",
		Env:               []string{"A=a", "B=b"},
		MinScale:          2,
		MaxScale:          2,
		ConcurrencyTarget: 100,
		ConcurrencyLimit:  100,
	}
	err := serving.CreateService(knclient, ksvcCfg)
	if err != nil {
		fmt.Println("create ksvc error:", err)
	}

	time.Sleep(time.Duration(10) * time.Second)

	ksvcCfg.Env = []string{"FUNCTION_ID=ewe-23e3-ereuwi-9495", "FUNCTION_LABEL=function=jim", "NAME=Haohao"}
	ksvcCfg.Image = "nginx:v2.0"
	err = serving.UpdateService(knclient, ksvcCfg)
	if err != nil {
		fmt.Println("update svc error:", err)
	}

	time.Sleep(time.Duration(10) * time.Second)

	err = serving.DeleteService(knclient, "jim")
	if err != nil {
		fmt.Println("delete svc error:", err)
	}
}
