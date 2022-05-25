#!/bin/bash

source './common.sh'

login

read -p "Enter devices name: " NAME
vault write identity/entity name=${NAME}.${DOMAIN} metadata=common_name=${NAME}.${DOMAIN}

export ID=$(vault read -format=json identity/entity/name/${NAME}.${DOMAIN} | jq -r .data.id)
export ACCESSOR=$(vault read -format=json /sys/auth | jq -r '.data."token/".accessor')
echo -e "${RED}Create entity alias"
vault write identity/entity-alias name=${NAME}.token canonical_id=${ID} mount_accessor=${ACCESSOR}
echo -e "${RED}Create token${NC}"
RES=$(vault write -format json auth/token/create/pki-client entity_alias=${NAME}.token)
TOKEN=$(echo $RES | jq -r '.auth.client_token')
echo -e "${RED}Generated token for $NAME${NC}: $TOKEN$"
mkdir -p client/${NAME}
echo $TOKEN > client/${NAME}/token
echo -e "${RED}Generating certificate${NC}"
export VAULT_TOKEN=${TOKEN}
vault write -format=json pki/issue/client common_name=${NAME}.${DOMAIN} ttl=12h | tee client/${NAME}/client.json
cat client/${NAME}/client.json  | jq -r .data.certificate > client/${NAME}/node.crt
cat client/${NAME}/client.json  | jq -r .data.private_key > client/${NAME}/node.key
cat client/${NAME}/client.json  | jq -r .data.issuing_ca > client/${NAME}/rootCA.crt
echo -e "${RED_BOLD}Find the files for ${NAME} in ./client/${NAME}${NC}"