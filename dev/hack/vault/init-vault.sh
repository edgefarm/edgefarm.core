#!/bin/bash

source './common.sh'

login

vault secrets enable pki
vault secrets tune -max-lease-ttl $((24*365*10))h pki
vault write pki/config/urls issuing_certificates=${FULL_URI}/cert crl_distribution_points=${FULL_URI}/v1/pki/crl ocsp_servers=${FULL_URI}/ocsp
vault write pki/root/generate/exported common_name=${COMMON_NAME} -format=json > ca.json
vault write pki/roles/server ext_key_usage=ServerAuth allowed_domains=${DOMAIN} allow_subdomains=true

echo -e "${RED}Writing policies${NC}"

# Allow clients to issue PKI client certs
vault policy write pki-client - << EOF
path "/pki/issue/client" {
    capabilities = [ "create", "update" ]
}
EOF

vault policy write pki-server - << EOF
path "/pki/issue/server" {
    capabilities = [ "create","update" ]
}
EOF

echo -e "${RED}Writing roles${NC}"

vault write pki/roles/client ext_key_usage=ClientAuth allowed_domains_template=true allowed_domains={{identity.entity.metadata.common_name}} allow_bare_domains=true
vault write pki/roles/server ext_key_usage=ServerAuth allowed_domains=${DOMAIN} allow_subdomains=true

echo -e "${RED}Writing token roles${NC}"
vault write auth/token/roles/pki-client allowed_policies=pki-client renewable=true token_explicit_max_ttl=24h token_no_default_policy=true allowed_entity_aliases="*.token"

echo -e "${RED}Enabling kubernets auth${NC}"
vault auth enable kubernetes
mkdir -p cluster
read -p "Enter cluster api server URI, e.g. 'https://cluster.example.com:6443': " KUBERNETES_URI

# export KUBERNETES_URI=$(kubectl cluster-info | grep "control plane" | awk -F " at " '{print $2}' | sed "s,\x1B\[[0-9;]*[a-zA-Z],,g")
export KUBERNETES_HOST_AND_PORT=$(echo ${KUBERNETES_URI} | awk -F "https://" '{print $2}')
openssl s_client -showcerts -connect ${KUBERNETES_HOST_AND_PORT} </dev/null 2>/dev/null |openssl x509 -outform PEM > cluster/ca.crt
echo -e "${RED}Setting kubernets auth config${NC}"
vault write auth/kubernetes/config kubernetes_host=${KUBERNETES_URI} kubernetes_ca_crt=cluster/ca.crt

echo -e "${RED}Create the binding role for the cloudcore service account${NC}"
vault write auth/kubernetes/role/cloudcore bound_service_account_names=cloudcore bound_service_account_namespaces=kubeedge token_policies=pki-server alias_name_source=serviceaccount_name token_no_default_policy=true

