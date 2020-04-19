package main

import (
	"fmt"
	"time"

	"github.com/qhh0205/kncrud/eventing"
)

// // //
// // type BrokerConfig struct {
// // 	Kind       string
// // 	Name       string
// // 	Namespace  string
// // 	APIVersion string
// // }

// // // Broker knative Broker common config
// // type Broker struct {
// // 	Name      string
// // 	Namespace string
// // 	Config    BrokerConfig
// // 	client    eventingv1beta1.BrokerInterface
// // }

// func main() {
// 	broker := eventing.NewBroker("default")
// 	broker.Name = "test-broker"
// 	broker.Config.Name = "config-br-defaults"
// 	broker.Config.Kind = "ConfigMap"
// 	broker.Config.Namespace = "knative-eventing"
// 	broker.Config.APIVersion = "v1"

// 	err := broker.Create()
// 	if err != nil {
// 		fmt.Println("create trigger error:", err)
// 	}

// 	time.Sleep(time.Duration(60) * time.Second)
// 	err = broker.Delete()
// 	if err != nil {
// 		fmt.Println("delete trigger error:", err)
// 	}
// }

// // Trigger knative trigger common config
// // Filter: 事件过滤器，和 cloudevent 的 context key-value 比较
// type Trigger struct {
// 	Name           string
// 	Namespace      string
// 	Broker         string
// 	Filter         map[string]string
// 	URL            string
// 	KRefKind       string
// 	KRefName       string
// 	KRefNamespace  string
// 	KRefAPIVersion string
// 	client         eventingv1beta1.TriggerInterface
// }

func main() {
	trigger := eventing.NewTrigger("default")
	trigger.Name = "test-trigger"
	trigger.Broker = "default"
	trigger.KRefKind = "Service"
	trigger.KRefName = "my-service"
	trigger.KRefAPIVersion = "serving.knative.dev/v1"
	trigger.KRefNamespace = "default"

	err := trigger.Create()
	if err != nil {
		fmt.Println("create trigger error:", err)
	}

	time.Sleep(time.Duration(60) * time.Second)
	trigger.URL = "/test"
	err = trigger.Update()
	if err != nil {
		fmt.Println("update trigger error:", err)
	}

	time.Sleep(time.Duration(60) * time.Second)
	err = trigger.Delete()
	if err != nil {
		fmt.Println("delete trigger error:", err)
	}
}

// // Function configmap data
// type Function struct {
// 	Name            string `json:"functionName"`
// 	CodeURL         string `json:"codeURL"`
// 	CodePath        string `json:"codePath"`
// 	CodeReleasePath string `json:"codeReleasePath"`
// 	Handler         string `json:"handler"` // handler format: module.handler
// 	ReleaseVersion  string `json:"releaseVersion"`
// }

// func main() {
// 	metaData := Function{
// 		"helloworld",
// 		"https://artifactory.bj.bcebos.com/index.py",
// 		"/home/work/code/function",
// 		"/home/work/code/release",
// 		"index.handler",
// 		"v1.0",
// 	}
// 	metaDataJSON, _ := json.Marshal(metaData)

// 	configmap := k8s.NewConfigmap("default")
// 	configmap.Name = "func1"
// 	configmap.Data = map[string]string{"meta.json": string(metaDataJSON)}
// 	configmap.Labels = map[string]string{"hello": "world"}

// 	err := configmap.Create()
// 	if err != nil {
// 		fmt.Println("create configmap error:", err)
// 	}

// 	time.Sleep(time.Duration(60) * time.Second)
// 	configmap.Labels = map[string]string{
// 		"hello":  "world",
// 		"funcid": "qwer",
// 	}
// 	err = configmap.Update()
// 	if err != nil {
// 		fmt.Println("update configmap error:", err)
// 	}

// 	time.Sleep(time.Duration(60) * time.Second)
// 	err = configmap.Delete()
// 	if err != nil {
// 		fmt.Println("delete configmap error:", err)
// 	}
// }

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/qhh0205/kncrud/k8s"
// 	// "icode.baidu.com/baidu/bpitsns/faas-controller/pkg/services/k8s"
// )

// // Function configmap data
// type Function struct {
// 	Name            string `json:"functionName"`
// 	CodeURL         string `json:"codeURL"`
// 	CodePath        string `json:"codePath"`
// 	CodeReleasePath string `json:"codeReleasePath"`
// 	Handler         string `json:"handler"` // handler format: module.handler
// 	ReleaseVersion  string `json:"releaseVersion"`
// }

// func main() {
// 	metaData := Function{
// 		"helloworld",
// 		"https://artifactory.bj.bcebos.com/index.py",
// 		"/home/work/code/function",
// 		"/home/work/code/release",
// 		"index.handler",
// 		"v1.0",
// 	}
// 	metaDataJSON, _ := json.Marshal(metaData)
// 	configmap := k8s.Configmap{
// 		Name:      "func1",
// 		Namespace: "default",
// 		Data:      map[string]string{"meta.json": string(metaDataJSON)},
// 		Labels:    map[string]string{"hello": "world"},
// 	}

// 	err := configmap.Create()
// 	if err != nil {
// 		fmt.Println("create configmap error:", err)
// 	}

// 	time.Sleep(time.Duration(60) * time.Second)
// 	configmap.Labels = map[string]string{
// 		"hello":  "world",
// 		"funcid": "qwer",
// 	}
// 	err = configmap.Update()
// 	if err != nil {
// 		fmt.Println("update configmap error:", err)
// 	}

// 	time.Sleep(time.Duration(60) * time.Second)
// 	err = configmap.Delete()
// 	if err != nil {
// 		fmt.Println("delete configmap error:", err)
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/qhh0205/kncrud/sources/pingsource"
// )

// // // PingSource knative pingsource common config
// // // URL 字段可选，可以是基于 Sink 的相对 URI，也可以是一个完整的 URL
// // type PingSource struct {
// // 	Name           string
// // 	Namespace      string
// // 	Schedule       string
// // 	JsonData       string
// // 	URL            string
// // 	KRefKind       string
// // 	KRefName       string
// // 	KRefNamespace  string
// // 	KRefAPIVersion string
// // }

// func main() {
// 	pingSource := pingsource.NewPingSource("default")
// 	pingSource.Name = "ping"
// 	pingSource.Schedule = "*/1 * * * *"
// 	pingSource.JsonData = `{"server":"127.0.0.1","server_port":1951}`
// 	pingSource.KRefKind = "Service"
// 	pingSource.KRefName = "event-display"
// 	pingSource.KRefAPIVersion = "serving.knative.dev/v1"

// 	fmt.Println(pingSource)
// 	err := pingSource.Create()
// 	if err != nil {
// 		fmt.Println("create pingsource error:", err)
// 	}

// 	time.Sleep(time.Duration(60) * time.Second)
// 	pingSource.URL = "/test"
// 	pingSource.Schedule = "*/2 * * * *"
// 	err = pingSource.Update()
// 	if err != nil {
// 		fmt.Println("update pingsource error:", err)
// 	}

// 	time.Sleep(time.Duration(60) * time.Second)
// 	err = pingSource.Delete()
// 	if err != nil {
// 		fmt.Println("delete pingsource error:", err)
// 	}
// }
