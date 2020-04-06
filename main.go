package main

import (
	"fmt"

	"github.com/qhh0205/kncrud/serving"
)

func main() {
	// ksvcCfg := serving.ServiceConfiguration{
	// 	Name:              "ccc",
	// 	Image:             "qhh0205/cloudevent:v2.0",
	// 	Envs:              map[string]string{"A": "a", "B": "b"},
	// 	MinScale:          2,
	// 	MaxScale:          2,
	// 	ConcurrencyTarget: 100,
	// 	ConcurrencyLimit:  100,
	// }
	// err := serving.CreateService(ksvcCfg, "default")
	// if err != nil {
	// 	fmt.Println("create ksvc error:", err)
	// }

	// time.Sleep(time.Duration(120) * time.Second)
	// ksvcCfg.Env = []string{"FUNCTION_ID=ewe-23e3-ereuwi-9495", "FUNCTION_LABEL=function=jim", "NAME=Haohao"}
	// ksvcCfg.Image = "qhh0205/cloudevent:v2.0"
	// ksvcCfg.Envs["A"] = "helloworld"
	// ksvcCfg.Envs["function"] = "xxxxxxxxxxxxxxxxxx"
	// ksvcCfg.MaxScale = 8
	// err := serving.UpdateService(ksvcCfg, "default")
	// if err != nil {
	// 	fmt.Println("update svc error:", err)
	// }

	// // time.Sleep(time.Duration(120) * time.Second)
	err := serving.DeleteService("ccc", "default")
	if err != nil {
		fmt.Println("delete svc error:", err)
	}
}
