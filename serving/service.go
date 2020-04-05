package serving

import (
	"flag"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	servinglib "knative.dev/client/pkg/serving"
	"knative.dev/client/pkg/util"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingclient "knative.dev/serving/pkg/client/clientset/versioned"
)

var knClientset *servingclient.Clientset

func init() {
	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "/Users/qianghaohao/.kube/config", "path to Kubernetes config file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	knClientset, err = servingclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

// ServiceConfiguration knative service common mete data
type ServiceConfiguration struct {
	// Direct field manipulation
	Name  string
	Image string
	Env   []string

	MinScale          int
	MaxScale          int
	ConcurrencyTarget int
	ConcurrencyLimit  int
}

// Apply Apply config
func (p *ServiceConfiguration) Apply(service *servingv1.Service) error {
	template := &service.Spec.Template
	err := servinglib.UpdateImage(template, p.Image)
	if err != nil {
		return err
	}

	envMap, err := util.MapFromArrayAllowingSingles(p.Env, "=")
	if err != nil {
		return err
	}
	envToRemove := util.ParseMinusSuffix(envMap)
	err = servinglib.UpdateEnvVars(template, envMap, envToRemove)
	if err != nil {
		return err
	}

	err = servinglib.UpdateMinScale(template, p.MinScale)
	if err != nil {
		return err
	}

	err = servinglib.UpdateMaxScale(template, p.MaxScale)
	if err != nil {
		return err
	}

	err = servinglib.UpdateConcurrencyTarget(template, p.ConcurrencyTarget)
	if err != nil {
		return err
	}

	err = servinglib.UpdateConcurrencyLimit(template, int64(p.ConcurrencyLimit))
	if err != nil {
		return err
	}

	return nil
}

// CreateService  create kantive service
func CreateService(serviceconfig ServiceConfiguration, namespace string) error {
	client := knClientset.ServingV1().Services(namespace)
	service, err := constructService(serviceconfig, namespace)
	if err != nil {
		return err
	}
	_, err = client.Create(service)
	if err != nil {
		return err
	}
	return nil
}

// UpdateService  delete kantive service
func UpdateService(serviceconfig ServiceConfiguration, namespace string) error {
	client := knClientset.ServingV1().Services(namespace)

	service, err := client.Get(serviceconfig.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	updatedService := service.DeepCopy()
	serviceconfig.Apply(updatedService)
	if err != nil {
		return err
	}
	_, err = client.Update(updatedService)
	if err != nil {
		return err
	}
	return nil
}

// DeleteService delete ksvc
func DeleteService(name string, namespace string) error {
	client := knClientset.ServingV1().Services(namespace)
	err := client.Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func constructService(serviceconfig ServiceConfiguration, namespace string) (*servingv1.Service,
	error) {
	service := servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceconfig.Name,
			Namespace: namespace,
		},
	}

	service.Spec.Template = servingv1.RevisionTemplateSpec{
		Spec:       servingv1.RevisionSpec{},
		ObjectMeta: metav1.ObjectMeta{},
	}
	service.Spec.Template.Spec.Containers = []corev1.Container{{}}

	err := serviceconfig.Apply(&service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}
