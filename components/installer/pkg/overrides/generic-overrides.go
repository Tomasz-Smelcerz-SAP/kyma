package overrides

import (
	"strings"

	"github.com/ghodss/yaml"
)

//Merges value into given map, introducing intermediate "nested" maps for every intermediate key.
func mergeIntoMap(keys []string, value string, dstMap map[string]interface{}) {
	currentKey := keys[0]
	//Last key points directly to string value
	if len(keys) == 1 {
		dstMap[currentKey] = value
		return
	}

	//All keys but the last one should point to a nested map
	nestedMap, ok := dstMap[currentKey].(map[string]interface{})

	if !ok {
		nestedMap = map[string]interface{}{}
		dstMap[currentKey] = nestedMap
	}

	mergeIntoMap(keys[1:], value, nestedMap)
}

func mapToYaml(sourceMap map[string]string) (string, error) {
	mergedMap := map[string]interface{}{}

	for key, value := range sourceMap {
		keys := strings.Split(key, ".")
		mergeIntoMap(keys, value, mergedMap)
	}

	res, err := yaml.Marshal(mergedMap)
	return string(res), err
}
