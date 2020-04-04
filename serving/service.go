package serving

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientservingv1 "knative.dev/client/pkg/serving/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// CreateService  create kantive service
func CreateService(client clientservingv1.KnServingClient, serviceconfig ServiceConfiguration) error {
	service, err := constructService(serviceconfig)
	if err != nil {
		return err
	}
	err = client.CreateService(service)
	if err != nil {
		return err
	}
	return nil
}

// UpdateService  delete kantive service
func UpdateService(client clientservingv1.KnServingClient, serviceconfig ServiceConfiguration) error {
	service, err := client.GetService(serviceconfig.Name)
	if err != nil {
		return err
	}
	updatedService := service.DeepCopy()
	serviceconfig.Apply(updatedService)
	if err != nil {
		return err
	}
	err = client.UpdateService(updatedService)
	if err != nil {
		return err
	}
	return nil
}

// DeleteService delete ksvc
func DeleteService(client clientservingv1.KnServingClient, name string) error {
	err := client.DeleteService(name, time.Duration(0))
	if err != nil {
		return err
	}
	return nil
}

func constructService(serviceconfig ServiceConfiguration) (*servingv1.Service,
	error) {
	service := servingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceconfig.Name,
			Namespace: serviceconfig.Namespace,
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
