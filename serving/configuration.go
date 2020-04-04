package serving

import (
	servinglib "knative.dev/client/pkg/serving"
	"knative.dev/client/pkg/util"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// ServiceConfiguration knative service common mete data
type ServiceConfiguration struct {
	// Direct field manipulation
	Name      string
	Namespace string
	Image     string
	Env       []string

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
