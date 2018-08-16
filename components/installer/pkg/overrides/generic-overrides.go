package overrides

import (
	"strings"

	"github.com/ghodss/yaml"
)

//Map of overrides. Values can be nested maps (of the same type) or strings
type OverridesMap map[string]interface{}

func unmarshallToNestedMap(value string) (OverridesMap, error) {
	target := OverridesMap{}

	err := yaml.Unmarshal([]byte(value), &target)

	return target, err
}

func overridesMapToYaml(oMap OverridesMap) (string, error) {
	res, err := yaml.Marshal(oMap)
	return string(res), err
}

func MergeMaps(base, newOverrides OverridesMap) {
	//TODO: Implement
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

func MarshallToYaml(oMap OverridesMap) string {
	//TODO: Implement
	return ""
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
