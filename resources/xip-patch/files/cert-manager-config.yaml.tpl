---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: kyma-ca-issuer
  namespace: istio-system
spec:
  ca:
    secretName: kyma-ca-key-pair
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: kyma-gateway-crt
  namespace: istio-system
spec:
  secretName: kyma-gateway-certs
  issuerRef:
    name: kyma-ca-issuer
    kind: Issuer
  commonName: "{{.Values.global.ingress.domainName}}"
  organization:
  - kymacia
  dnsNames:
  - "*.{{.Values.global.ingress.domainName}}"

