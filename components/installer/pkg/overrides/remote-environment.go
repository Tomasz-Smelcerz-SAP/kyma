package overrides

const hmcDefault = `
deployment:
  args:
    sourceType: marketing
service:
  externalapi:
    nodePort: 32001
`

const ecDefault = `
deployment:
  args:
    sourceType: commerce
service:
  externalapi:
    nodePort: 32000
`

// GetHmcDefaultOverrides returns values overrides for hmc default remote environment
func GetHmcDefaultOverrides() (OverridesMap, error) {
	oMap, err := unmarshallToNestedMap(hmcDefault)
	if err != nil {
		return nil, err
	}
	return oMap, nil
}

// GetEcDefaultOverrides returns values overrides for ec default remote environment
func GetEcDefaultOverrides() (OverridesMap, error) {
	oMap, err := unmarshallToNestedMap(ecDefault)
	if err != nil {
		return nil, err
	}
	return oMap, nil
}
