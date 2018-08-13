package overrides

import (
	core "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const label = "installer: overrides"

// ReaderInterface exposes functions
type ReaderInterface interface {
	GetFullConfigMap() (map[string]string, error)
}

type reader struct {
	configmaps v1.ConfigMapInterface
	secrets    v1.SecretInterface
}

// NewReader returns a ready to use configmapClient
func NewReader(namespace string, client v1.CoreV1Client) (ReaderInterface, error) {
	r := &reader{
		configmaps: client.ConfigMaps(namespace),
		secrets:    client.Secrets(namespace),
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

func (r reader) getLabeledConfigMaps() ([]core.ConfigMap, error) {
	listOpts := meta_v1.ListOptions{
		LabelSelector: label,
	}
	configmaps, err := r.configmaps.List(listOpts)
	if err != nil {
		return nil, err
	}
	return configmaps.Items, nil
}

func (r reader) getLabeledSecrets() ([]core.Secret, error) {
	listOpts := meta_v1.ListOptions{
		LabelSelector: label,
	}
	secrets, err := r.secrets.List(listOpts)
	if err != nil {
		return nil, err
	}
	return secrets.Items, nil
}
