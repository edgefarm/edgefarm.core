#!/bin/bash
RED='\033[0;31m'
RED_BOLD='\033[1;31m'
NC='\033[0m' # No Color

NS=vault-system

login() {
    FULL_URI="http://192-168-1-42.nip.io:38200"
    
    # must only be exactly one!!
    POD=$(kubectl get pods -n ${NS} --no-headers -o custom-columns=":metadata.name")
    ROOT_TOKEN=$(kubectl logs -n ${NS} ${POD} | grep "Root Token" | awk -F " " '{print $3}')
    
    export VAULT_TOKEN=${ROOT_TOKEN}
    # read -p "Enter hosts full URI, e.g. https://vault.example.com" FULL_URI
    
    # exctract the common name from the FULL_URI, e.g. https://vault.example.com -> vault.example.com
    COMMON_NAME=$(echo ${FULL_URI} | awk -F "/" '{print $3}' | awk -F ":" '{print $1}')
    # extract the domain from the COMMON_NAME, e.g. vault.example.com -> example.com
    DOMAIN=$(echo ${COMMON_NAME} | awk -F "."  '{print $2"."$3}')
    echo $DOMAIN
    export VAULT_ADDR=${FULL_URI}
    
    vault login ${ROOT_TOKEN}
}