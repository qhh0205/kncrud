package eventing

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "knative.dev/eventing/pkg/apis/eventing/v1beta1"
	eventingv1beta1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/eventing/v1beta1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// Trigger knative trigger common config
// Filter: 事件过滤器，和 cloudevent 的 context key-value 比较
type Trigger struct {
	Name           string
	Namespace      string
	Broker         string
	Filter         map[string]string
	URL            string
	KRefKind       string
	KRefName       string
	KRefNamespace  string
	KRefAPIVersion string
	client         eventingv1beta1.TriggerInterface
}

// NewTrigger new Trigger object
func NewTrigger(namespace string) *Trigger {
	return &Trigger{
		client:    knClientset.EventingV1beta1().Triggers(namespace),
		Namespace: namespace,
	}
}

// Create create Trigger
func (tg *Trigger) Create() error {
	trigger, err := constructTrigger(tg)
	if err != nil {
		return err
	}
	_, err = tg.client.Create(trigger)
	if err != nil {
		return err
	}
	return nil
}

// Update update Trigger
func (tg *Trigger) Update() error {
	uri, err := apis.ParseURL(tg.URL)
	if err != nil {
		return err
	}
	retries := 0
	nrRetries := 3
	for {
		trigger, err := tg.client.Get(tg.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updatedTrigger := trigger.DeepCopy()

		updatedTrigger.Spec.Broker = tg.Broker
		updatedTrigger.Spec.Filter.Attributes = tg.Filter
		updatedTrigger.Spec.Subscriber = duckv1.Destination{
			Ref: &duckv1.KReference{
				Kind:       tg.KRefKind,
				Name:       tg.KRefName,
				Namespace:  tg.KRefNamespace,
				APIVersion: tg.KRefAPIVersion,
			},
			URI: uri,
		}
		_, err = tg.client.Update(updatedTrigger)
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

// Delete delete Trigger
func (tg *Trigger) Delete() error {
	deletePolicy := metav1.DeletePropagationBackground
	err := tg.client.Delete(tg.Name, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		return err
	}
	return nil
}

func constructTrigger(tg *Trigger) (*v1beta1.Trigger, error) {
	uri, err := apis.ParseURL(tg.URL)
	if err != nil {
		return nil, err
	}
	return &v1beta1.Trigger{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tg.Name,
			Namespace: tg.Namespace,
		},
		Spec: v1beta1.TriggerSpec{
			Broker: tg.Broker,
			Filter: &v1beta1.TriggerFilter{
				Attributes: tg.Filter,
			},
			Subscriber: duckv1.Destination{
				Ref: &duckv1.KReference{
					Kind:       tg.KRefKind,
					Name:       tg.KRefName,
					Namespace:  tg.KRefNamespace,
					APIVersion: tg.KRefAPIVersion,
				},
				URI: uri,
			},
		},
	}, nil
}
