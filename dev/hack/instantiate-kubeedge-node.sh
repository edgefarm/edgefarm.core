#!/bin/bash
CLOUDCORE_ADDRESS_DOTTED=${1}
CLOUDCORE_ADDRESS=$(echo ${CLOUDCORE_ADDRESS_DOTTED} | sed -e 's/\./\-/g' -e 's/-nip-io/.nip.io/g')
NODE_NAME=${2}
docker run -d --rm --env CLOUDCORE_ADDRESS=${CLOUDCORE_ADDRESS} --env NODE_NAME=${NODE_NAME} --name ${NODE_NAME} \
-v ${HOME}/.devspace/ca/kubeedge/rootCA.pem:/etc/kubeedge/certs/rootCa.pem \
-v $(pwd)/hack/kubeedge-node/certs/node.pem:/etc/kubeedge/certs/node.pem \
-v $(pwd)/hack/kubeedge-node/certs/node.key:/etc/kubeedge/certs/node.key \
--privileged edgefarm/virtual-device:latest
