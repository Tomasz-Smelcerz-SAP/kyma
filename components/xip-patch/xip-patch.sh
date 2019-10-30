#!/usr/bin/env bash

set -o errexit

# # # # # # # # # # # # # # # # #
# VARs coming from environment: #
#                               #
# EXTERNAL_PUBLIC_IP            #
# INGRESSGATEWAY_SERVICE_NAME   #
# GLOBAL_DOMAIN                 #
# GLOBAL_TLS_CERT               #
# GLOBAL_TLS_KEY                #
# INGRESS_DOMAIN                #
# INGRESS_TLS_CERT              #
# INGRESS_TLS_KEY               #
# # # # # # # # # # # # # # # # #

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $CURRENT_DIR/utils.sh

generateXipDomain() {

    if [ -z "${EXTERNAL_PUBLIC_IP}" ]; then

        local namespace="istio-system"

        if [ -z "${INGRESSGATEWAY_SERVICE_NAME}" ]; then
            INGRESSGATEWAY_SERVICE_NAME=istio-ingressgateway
        fi

        EXTERNAL_PUBLIC_IP=$(getLoadBalancerIP "${INGRESSGATEWAY_SERVICE_NAME}" "${namespace}")

        if [[ "$?" != 0 ]]; then
            echo "External public IP not found"
            exit 1
        fi
    fi

    echo "${EXTERNAL_PUBLIC_IP}.xip.io"

}

generateCerts() {
    TEMP=$(mktemp /tmp/cert-file.XXXXXXXX)
    sed 's/{{.Values.global.ingress.domainName}}/'$INGRESS_DOMAIN'/' /etc/cert-config/config.yaml.tpl > ${TEMP}

    echo DEBUG:
    cat ${TEMP}
    kubectl create -f ${TEMP}
    rm ${TEMP}
}

createOverridesConfigMap() {
    if [ -z "$(kubectl get configmap -n kyma-installer net-global-overrides --ignore-not-found)" ]; then
        kubectl create configmap net-global-overrides \
            --from-literal global.ingress.domainName="$INGRESS_DOMAIN" \
            --from-literal global.ingress.tlsCrt="$INGRESS_TLS_CERT" \
            --from-literal global.ingress.tlsKey="$INGRESS_TLS_KEY" \
            -n kyma-installer
    fi
    kubectl label configmap net-global-overrides --overwrite installer=overrides -n kyma-installer
    kubectl label configmap net-global-overrides --overwrite kyma-project.io/installation="" -n kyma-installer
}

patchTlsCrtSecret() {
    TLS_CERT_YAML=$(cat << EOF
---
data:
  tls.crt: "${INGRESS_TLS_CERT}"
EOF
    )
    set +e
    local msg
    local status
    msg=$(kubectl patch secret ingress-tls-cert --patch "${TLS_CERT_YAML}" -n kyma-system 2>&1)
    status=$?
    set -e
    if [[ $status -ne 0 ]] && [[ ! "$msg" == *"not patched"* ]]; then
        echo "$msg"
        exit $status
    fi
}

#This does not exist on install (takes fallback value), but it exists on update!
INGRESS_TLS_CERT="${INGRESS_TLS_CERT:-$GLOBAL_TLS_CERT}"
INGRESS_TLS_KEY="${INGRESS_TLS_KEY:-$GLOBAL_TLS_KEY}"
INGRESS_DOMAIN="${INGRESS_DOMAIN:-$GLOBAL_DOMAIN}"

if [ -n "${INGRESS_TLS_CERT}" ] && [ -z "${INGRESS_DOMAIN}" ]; then
    echo "Certificate provided, but domain is missing!"
    exit 1
fi

if [ -z "${INGRESS_DOMAIN}" ] ; then
    INGRESS_DOMAIN=$(generateXipDomain)
fi

if [ -z "${INGRESS_TLS_CERT}" ] ; then
    generateCerts
fi

createOverridesConfigMap

patchTlsCrtSecret

