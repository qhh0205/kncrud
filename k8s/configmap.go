package k8s

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var clientset *kubernetes.Clientset

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

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

// Configmap configmap common meta data
type Configmap struct {
	Name      string
	Namespace string
	Data      map[string]string
	Labels    map[string]string
	client    corev1client.ConfigMapInterface
}

// NewConfigmap reutrn a Configmap object with client
func NewConfigmap(namespace string) *Configmap {
	return &Configmap{
		client:    clientset.CoreV1().ConfigMaps(namespace),
		Namespace: namespace,
	}
}

// Create create configmap
func (cm *Configmap) Create() error {
	configmap := constructConfigmap(cm)
	// Note: 从 0.18.x 开始需要三个参数
	// https://github.com/kubernetes/client-go/blob/v0.18.0/kubernetes/typed/core/v1/configmap.go
	_, err := cm.client.Create(configmap)
	if err != nil {
		return err
	}
	return nil
}

// Update update configmap
func (cm *Configmap) Update() error {
	retries := 0
	nrRetries := 3
	for {
		configmap, err := cm.client.Get(cm.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updatedConfigmap := configmap.DeepCopy()
		updatedConfigmap.Data = cm.Data
		updatedConfigmap.Labels = cm.Labels
		_, err = cm.client.Update(updatedConfigmap)
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

// Delete delete configmap
func (cm *Configmap) Delete() error {
	deletePolicy := metav1.DeletePropagationBackground
	err := cm.client.Delete(cm.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}
	return nil
}

// constructConfigmap construct configmap with common data
func constructConfigmap(cm *Configmap) *apiv1.ConfigMap {
	return &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cm.Name,
			Namespace: cm.Namespace,
			Labels:    cm.Labels,
		},
		Data: cm.Data,
	}
}
