package overrides

import (
	"bytes"
	"text/template"

	"github.com/kyma-project/kyma/components/installer/pkg/config"
)

const globalsTplStr = `
global:
  tlsCrt: "{{.ClusterTLSCert}}"
  tlsKey: "{{.ClusterTLSKey}}"
  isLocalEnv: {{.IsLocalInstallation}}
  domainName: "{{.Domain}}"
  remoteEnvCa: "{{.RemoteEnvCa}}"
  remoteEnvCaKey: "{{.RemoteEnvCaKey}}"
  istio:
    tls:
      secretName: "istio-ingress-certs"
  etcdBackupABS:
    containerName: "{{.EtcdBackupABSContainerName}}"
`

// GetGlobalOverrides .
func GetGlobalOverrides(installationData *config.InstallationData) (OverridesMap, error) {

	tmpl, err := template.New("").Parse(globalsTplStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, installationData)
	if err != nil {
		return nil, err
	}

	oMap, err := unmarshallToNestedMap(buf.String())
	if err != nil {
		return nil, err
	}
	return oMap, nil
}
