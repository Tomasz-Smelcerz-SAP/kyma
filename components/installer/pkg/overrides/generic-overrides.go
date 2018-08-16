package overrides

import (
	"strings"

	"github.com/ghodss/yaml"
)

//Map of overrides. Values can be nested maps (of the same type) or strings
type OverridesMap map[string]interface{}

func ToMap(value string) (OverridesMap, error) {
	target := OverridesMap{}

	err := yaml.Unmarshal([]byte(value), &target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func ToYaml(oMap OverridesMap) (string, error) {
	res, err := yaml.Marshal(oMap)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
func MergeMaps(baseMap, newOnes OverridesMap) {
	//TODO: Implement
	//for key, value := range newOnes
}

//Merges value into given map, introducing intermediate "nested" maps for every intermediate key.
func mergeIntoMap(keys []string, value string, dstMap OverridesMap) {
	currentKey := keys[0]
	//Last key points directly to string value
	if len(keys) == 1 {
		dstMap[currentKey] = value
		return
	}

	//All keys but the last one should point to a nested map
	nestedMap, ok := dstMap[currentKey].(OverridesMap)

	if !ok {
		nestedMap = OverridesMap{}
		dstMap[currentKey] = nestedMap
	}

	mergeIntoMap(keys[1:], value, nestedMap)
}

//Used to convert external "flat" overrides into OverridesMap.
func flatMapToOverridesMap(sourceMap map[string]string) OverridesMap {
	mergedMap := OverridesMap{}
	if len(sourceMap) == 0 {
		return mergedMap
	}

	for key, value := range sourceMap {
		keys := strings.Split(key, ".")
		mergeIntoMap(keys, value, mergedMap)
	}

	return mergedMap
}
