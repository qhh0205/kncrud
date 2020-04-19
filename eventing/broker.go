package eventing

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "knative.dev/eventing/pkg/apis/eventing/v1beta1"
	eventingv1beta1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/eventing/v1beta1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// BrokerConfig Broker 配置，指向 Configmap
type BrokerConfig struct {
	Kind       string
	Name       string
	Namespace  string
	APIVersion string
}

// Broker knative Broker common config
// TODO; 先不考虑死信队列相关配置
type Broker struct {
	Name      string
	Namespace string
	Config    BrokerConfig
	client    eventingv1beta1.BrokerInterface
}

// NewBroker new Broker object
func NewBroker(namespace string) *Broker {
	return &Broker{
		client:    knClientset.EventingV1beta1().Brokers(namespace),
		Namespace: namespace,
	}
}

// Create create Broker
func (br *Broker) Create() error {
	broker := constructBroker(br)
	_, err := br.client.Create(broker)
	if err != nil {
		return err
	}
	return nil
}

// Update update Broker
func (br *Broker) Update() error {
	retries := 0
	nrRetries := 3
	for {
		broker, err := br.client.Get(br.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updatedBroker := broker.DeepCopy()
		updatedBroker.Spec.Config = &duckv1.KReference{
			Kind:       br.Config.Kind,
			Name:       br.Config.Name,
			Namespace:  br.Config.Namespace,
			APIVersion: br.Config.APIVersion,
		}
		_, err = br.client.Update(updatedBroker)
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

// Delete delete Broker
func (br *Broker) Delete() error {
	deletePolicy := metav1.DeletePropagationBackground
	err := br.client.Delete(br.Name, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		return err
	}
	return nil
}

func constructBroker(br *Broker) *v1beta1.Broker {
	return &v1beta1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:      br.Name,
			Namespace: br.Namespace,
		},
		Spec: v1beta1.BrokerSpec{
			Config: &duckv1.KReference{
				Kind:       br.Config.Kind,
				Name:       br.Config.Name,
				Namespace:  br.Config.Namespace,
				APIVersion: br.Config.APIVersion,
			},
		},
	}
}
