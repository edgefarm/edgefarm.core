#!/bin/bash
#
# Execute this script for provisioning edgefarm certs to a device.
#
# To be called with arguments: <device ip address>
#
set -e
# Requires generated certificates:
CAROOT=${HOME}/.devspace/ca/kubeedge
rootCa=${CAROOT}/rootCA.pem
nodePublicKey=$(pwd)/dev/hack/kubeedge-node/certs/node.pem
nodePrivateKey=$(pwd)/dev/hack/kubeedge-node/certs/node.key
#

if [ "$#" -ne 1 ] ; then
    echo "Usage: ${0} <device ip address>"
    exit 1
fi

deviceIP=${1}

if [ ! -f "$rootCa" ]; then
    echo "File $rootCa does not exist."
    echo "Creating self-signed CA kubeedge"
    CAROOT=${CAROOT} mkcert 2>/dev/null
fi

if [ ! -f "$nodePublicKey" ] ||  [ ! -f "$nodePrivateKey" ]; then
    echo "Missing client certificate."
    echo "Generate kubeedge client certs."
    devspace run mkcert.create-client-cert kubeedge \
    -cert-file ${nodePublicKey} \
    -key-file ${nodePrivateKey} *.nip.io *.edgefarm.local
fi

scp ${rootCa} root@${deviceIP}:/etc/kubeedge/certs/
scp ${nodePublicKey} root@${deviceIP}:/etc/kubeedge/certs/
scp ${nodePrivateKey} root@${deviceIP}:/etc/kubeedge/certs/
