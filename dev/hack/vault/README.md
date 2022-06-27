Hack vault
==========

There are three scripts within this folder

`init-vault.sh` configures the vault server to use the PKI and sets all articats (policies and roles) for pki-client and pki-server.
pki-client shall be used for kubeedge devices to receive their certs.
pki-server shall be used to the kubeedge cloudcore to receive its certs.

`cloudcore.sh` sets all artifacts that cloudcore can retrieve certificates.

`newEdge.sh` creates a new token and certificate for an edge device. All files generated from `client/<name>` can be copied on to the edge device.