package serving

import (
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"knative.dev/pkg/ptr"
	"knative.dev/serving/pkg/apis/autoscaling"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingclient "knative.dev/serving/pkg/client/clientset/versioned"
)

var knClientset *servingclient.Clientset

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
	Envs  map[string]string

	MinScale          int
	MaxScale          int
	ConcurrencyTarget int
	ConcurrencyLimit  int
}

// CreateService  create kantive service
func CreateService(serviceconfig ServiceConfiguration, namespace string) error {
	client := knClientset.ServingV1().Services(namespace)
	service := constructService(serviceconfig, namespace)
	_, err := client.Create(service)
	if err != nil {
		return err
	}
	return nil
}

// UpdateService  delete kantive service
func UpdateService(serviceconfig ServiceConfiguration, namespace string) error {
	client := knClientset.ServingV1().Services(namespace)
	retries := 0
	nrRetries := 3
	for {
		service, err := client.Get(serviceconfig.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updatedService := service.DeepCopy()
		envVars := map2ContainerEnvVar(serviceconfig.Envs)
		updatedTemplate := &updatedService.Spec.Template
		updatedTemplate.Spec.Containers[0].Image = serviceconfig.Image
		updatedTemplate.Spec.Containers[0].Env = envVars
		updatedTemplate.Annotations[autoscaling.MaxScaleAnnotationKey] = strconv.Itoa(serviceconfig.MaxScale)
		updatedTemplate.Annotations[autoscaling.MinScaleAnnotationKey] = strconv.Itoa(serviceconfig.MinScale)
		updatedTemplate.Annotations[autoscaling.TargetAnnotationKey] = strconv.Itoa(serviceconfig.ConcurrencyTarget)
		updatedTemplate.Spec.ContainerConcurrency = ptr.Int64(int64(serviceconfig.ConcurrencyLimit))

		_, err = client.Update(updatedService)
		if err != nil {
			if apierrors.IsConflict(err) && retries < nrRetries {
				retries++
				// Wait a second before doing the retry
				time.Sleep(time.Second)
				continue
			}
			return errors.Wrap(err, fmt.Sprintf("giving up after %d retries", nrRetries))
		}
		return nil
	}
}

// DeleteService delete ksvc
func DeleteService(name string, namespace string) error {
	client := knClientset.ServingV1().Services(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := client.Delete(name, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		return err
	}
	return nil
}

func constructService(serviceconfig ServiceConfiguration, namespace string) *servingv1.Service {
	envVars := map2ContainerEnvVar(serviceconfig.Envs)
	return &servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceconfig.Name,
			Namespace: namespace,
		},
		Spec: servingv1.ServiceSpec{
			ConfigurationSpec: servingv1.ConfigurationSpec{
				Template: servingv1.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							autoscaling.MaxScaleAnnotationKey: strconv.Itoa(serviceconfig.MaxScale),
							autoscaling.MinScaleAnnotationKey: strconv.Itoa(serviceconfig.MinScale),
							autoscaling.TargetAnnotationKey:   strconv.Itoa(serviceconfig.ConcurrencyTarget),
						},
					},
					Spec: servingv1.RevisionSpec{
						ContainerConcurrency: ptr.Int64(int64(serviceconfig.ConcurrencyLimit)),
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Image: serviceconfig.Image,
									Env:   envVars,
								},
							},
						},
					},
				},
			},
		},
	}
}

// map2ContainerEnvVar map 转换为容器环境变量列表 key: Nmae, value: Value
func map2ContainerEnvVar(m map[string]string) []corev1.EnvVar {
	var envVars []corev1.EnvVar
	for k, v := range m {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}
	return envVars
}
