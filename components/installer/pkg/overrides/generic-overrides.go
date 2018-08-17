package overrides

import (
	"strings"

	"github.com/ghodss/yaml"
)

//OverridesMap is a map of overrides. Values in the map can be nested maps (of the same type) or strings
type OverridesMap map[string]interface{}

//ToMap converts yaml to OverridesMap. Supports only map-like yamls (no lists!)
func ToMap(value string) (OverridesMap, error) {
	target := OverridesMap{}

	err := yaml.Unmarshal([]byte(value), &target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

//ToYaml converts OverridesMap to yaml
func ToYaml(oMap OverridesMap) (string, error) {
	if len(oMap) == 0 {
		return "", nil
	}

	res, err := yaml.Marshal(oMap)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func FlattenMap(oMap OverridesMap) map[string]string {
	res := map[string]string{}
	_flattenMap(oMap, "", res)
	return res
}

//Used to convert external "flat" overrides into OverridesMap.
func UnflattenMap(sourceMap map[string]string) OverridesMap {
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

func mergeInto(baseMap map[string]interface{}, key string, newVal interface{}) {

	baseVal := baseMap[key]

	bVal, baseIsMap := baseVal.(map[string]interface{})
	nVal, newIsMap := newVal.(map[string]interface{})

	if baseIsMap && newIsMap {
		//Two maps case! Mutual Reccursion here :)
		MergeMaps(bVal, nVal)
	} else {
		//All other cases
		baseMap[key] = newVal
	}
}

//MergeMaps merges all values from newOnes map into baseMap, overwriting final keys (string values) if both maps contain such entries
func MergeMaps(baseMap, newOnes OverridesMap) {
	for key, newVal := range newOnes {
		_, baseContains := baseMap[key]
		if baseContains {
			//baseMap contain the entry.
			mergeInto(baseMap, key, newVal)
		} else {
			//baseMap does not contain such entry. Just use newVal and we're done.
			baseMap[key] = newVal
		}
	}
}

// Flattens given OverridesMap. The keys in result map will contain all intermediate keys joined with dots, e.g.: "istio.ingress.service.gateway: xyz"
func _flattenMap(oMap OverridesMap, keys string, result map[string]string) {

	var prefix string

	if len(keys) == 0 {
		prefix = ""
	} else {
		prefix = keys + "."
	}

	for key, value := range oMap {

		aString, isString := value.(string)
		if isString {
			result[prefix+key] = aString
		} else {
			//Nested map!
			nestedMap := value.(map[string]interface{})
			_flattenMap(nestedMap, prefix+key, result)
		}
	}
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
