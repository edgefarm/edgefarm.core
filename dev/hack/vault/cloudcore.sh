#!/bin/bash

source './common.sh'

login

read -p "Enter hostname of cloudcore, e.g. cloudcore.example.com: " CLOUDCORE_HOSTNAME

vault write identity/entity name=${CLOUDCORE_HOSTNAME}
export ID=$(vault read -format=json identity/entity/name/${CLOUDCORE_HOSTNAME} | jq -r .data.id)
export ACCESSOR=$(vault auth list -format=json|jq -r '.["kubernetes/"].accessor')
vault write identity/entity-alias canonical_id=$ID mount_accessor=$ACCESSOR name=kubeedge/cloudcore
