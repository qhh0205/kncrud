package pingsource

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	v1alpha2 "knative.dev/eventing/pkg/apis/sources/v1alpha2"
	eventingclient "knative.dev/eventing/pkg/client/clientset/versioned"
	sourcesv1alpha2 "knative.dev/eventing/pkg/client/clientset/versioned/typed/sources/v1alpha2"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
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

// PingSource knative pingsource common config
// URL 字段可选，可以是基于 Sink 的相对 URI，也可以是一个完整的 URL
type PingSource struct {
	Name           string
	Namespace      string
	Schedule       string
	JsonData       string
	URL            string
	KRefKind       string
	KRefName       string
	KRefNamespace  string
	KRefAPIVersion string
	client         sourcesv1alpha2.PingSourceInterface
}

// NewPingSource new PingSource object
func NewPingSource(namespace string) *PingSource {
	return &PingSource{
		client:    knClientset.SourcesV1alpha2().PingSources(namespace),
		Namespace: namespace,
	}
}

// Create create PingSource
func (ps *PingSource) Create() error {
	pingSource, err := constructPingSource(ps)
	if err != nil {
		return err
	}
	_, err = ps.client.Create(pingSource)
	if err != nil {
		return err
	}
	return nil
}

// Update update PingSource
func (ps *PingSource) Update() error {
	uri, err := apis.ParseURL(ps.URL)
	if err != nil {
		return err
	}
	retries := 0
	nrRetries := 3
	for {
		pingSource, err := ps.client.Get(ps.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updatedPingSource := pingSource.DeepCopy()
		updatedPingSource.Spec.Schedule = ps.Schedule
		updatedPingSource.Spec.JsonData = ps.JsonData
		updatedPingSource.Spec.Sink = duckv1.Destination{
			Ref: &duckv1.KReference{
				Kind:       ps.KRefKind,
				Name:       ps.KRefName,
				Namespace:  ps.KRefNamespace,
				APIVersion: ps.KRefAPIVersion,
			},
			URI: uri,
		}
		_, err = ps.client.Update(updatedPingSource)
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

// Delete delete PingSource
func (ps *PingSource) Delete() error {
	deletePolicy := metav1.DeletePropagationBackground
	err := ps.client.Delete(ps.Name, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		return err
	}
	return nil
}

func constructPingSource(ps *PingSource) (*v1alpha2.PingSource, error) {
	uri, err := apis.ParseURL(ps.URL)
	if err != nil {
		return nil, err
	}
	return &v1alpha2.PingSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ps.Name,
			Namespace: ps.Namespace,
		},
		Spec: v1alpha2.PingSourceSpec{
			SourceSpec: duckv1.SourceSpec{
				Sink: duckv1.Destination{
					Ref: &duckv1.KReference{
						Kind:       ps.KRefKind,
						Name:       ps.KRefName,
						Namespace:  ps.KRefNamespace,
						APIVersion: ps.KRefAPIVersion,
					},
					URI: uri,
				},
			},
			Schedule: ps.Schedule,
			JsonData: ps.JsonData,
		},
	}, nil
}
