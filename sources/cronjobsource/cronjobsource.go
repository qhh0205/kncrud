package cronjobsource

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
	v1alpha1 "knative.dev/eventing/pkg/apis/sources/v1alpha1"
	eventingclient "knative.dev/eventing/pkg/client/clientset/versioned"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
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

// CronJobSourceConfig knative cronjobsource common config
type CronJobSourceConfig struct {
	Name                 string
	Schedule             string
	Data                 string
	DeprecatedKind       string
	DeprecatedName       string
	DeprecatedNamespace  string
	DeprecatedAPIVersion string
}

// CreateCronJobSource  create cronjobsource
func CreateCronJobSource(cronconfig CronJobSourceConfig, namespace string) error {
	client := knClientset.SourcesV1alpha1().CronJobSources(namespace)
	cronjobsource := constructCronJobSource(cronconfig, namespace)
	_, err := client.Create(cronjobsource)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCronJobSource update conjobsource
func UpdateCronJobSource(cronconfig CronJobSourceConfig, namespace string) error {
	client := knClientset.SourcesV1alpha1().CronJobSources(namespace)
	retries := 0
	nrRetries := 3
	for {
		cronjobsource, err := client.Get(cronconfig.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		updatedCronjobsource := cronjobsource.DeepCopy()
		updatedCronjobsource.Spec.Schedule = cronconfig.Schedule
		updatedCronjobsource.Spec.Data = cronconfig.Data
		updatedCronjobsource.Spec.Sink = &duckv1beta1.Destination{
			DeprecatedKind:       cronconfig.DeprecatedKind,
			DeprecatedName:       cronconfig.DeprecatedName,
			DeprecatedNamespace:  cronconfig.DeprecatedNamespace,
			DeprecatedAPIVersion: cronconfig.DeprecatedAPIVersion,
		}
		_, err = client.Update(updatedCronjobsource)
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

// DeleteCronJobSource delete cronjobsource
func DeleteCronJobSource(name string, namespace string) error {
	client := knClientset.SourcesV1alpha1().CronJobSources(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := client.Delete(name, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		return err
	}
	return nil
}

// knative 0.11 之后将 CronJobSource 改为了 PingSource...
func constructCronJobSource(cronconfig CronJobSourceConfig, namespace string) *v1alpha1.CronJobSource {
	return &v1alpha1.CronJobSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cronconfig.Name,
			Namespace: namespace,
		},
		Spec: v1alpha1.CronJobSourceSpec{
			Schedule: cronconfig.Schedule,
			Data:     cronconfig.Data,
			Sink: &duckv1beta1.Destination{
				DeprecatedKind:       cronconfig.DeprecatedKind,
				DeprecatedName:       cronconfig.DeprecatedName,
				DeprecatedNamespace:  cronconfig.DeprecatedNamespace,
				DeprecatedAPIVersion: cronconfig.DeprecatedAPIVersion,
			},
		},
	}
}
