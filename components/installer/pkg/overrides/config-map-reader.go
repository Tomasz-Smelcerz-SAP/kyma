package overrides

import (
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/informers"
	listers "k8s.io/client-go/listers/core/v1"
)

const (
	key   = "installer"
	value = "overrides"
)

// ReaderInterface exposes functions
type ReaderInterface interface {
	GetFullConfigMap() (map[string]string, error)
}

type reader struct {
	configmaps listers.ConfigMapLister
	secrets    listers.SecretLister
}

// NewReader returns a ready to use configmapClient
func NewReader(namespace string, kubeInformerFactory informers.SharedInformerFactory) (ReaderInterface, error) {

	configmapLister := kubeInformerFactory.Core().V1().ConfigMaps().Lister()
	secretLister := kubeInformerFactory.Core().V1().Secrets().Lister()

	r := &reader{
		configmaps: configmapLister,
		secrets:    secretLister,
	}

	return r, nil
}

func (r reader) GetFullConfigMap() (map[string]string, error) {

	var combined map[string]string

	configmaps, err := r.getLabeledConfigMaps()
	if err != nil {
		return nil, err
	}

	secrets, err := r.getLabeledSecrets()
	if err != nil {
		return nil, err
	}

	for _, cMap := range configmaps {
		for key, val := range cMap.Data {
			combined[key] = val
		}
	}

	for _, sec := range secrets {
		for key, val := range sec.StringData {
			combined[key] = val
		}
	}

	return combined, nil
}

func (r reader) getLabeledConfigMaps() ([]*core.ConfigMap, error) {

	selector, err := getSelector()
	if err != nil {
		return nil, err
	}

	configmaps, err := r.configmaps.List(selector)
	if err != nil {
		return nil, err
	}

	return configmaps, nil
}

func (r reader) getLabeledSecrets() ([]*core.Secret, error) {

	selector, err := getSelector()
	if err != nil {
		return nil, err
	}

	secrets, err := r.secrets.List(selector)
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func getSelector() (labels.Selector, error) {

	req, err := labels.NewRequirement(key, selection.Equals, []string{value})
	if err != nil {
		return nil, err
	}
	selector := labels.NewSelector().Add(*req)

	return selector, nil
}
