---
apiVersion: certmanager.k8s.io/v1alpha1
kind: ClusterIssuer
metadata:
  name: selfsigning-issuer
  namespace: istio-system
spec:
  selfSigned: {}
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: kyma-gateway-crt
  namespace: istio-system
spec:
  secretName: kyma-gateway-certs
  commonName: "{{.Values.global.ingress.domainName}}"
  dnsNames:
  - "*.{{.Values.global.ingress.domainName}}"
  isCA: true
  issuerRef:
    name: selfsigning-issuer
    kind: ClusterIssuer

