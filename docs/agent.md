Agents
======

Agents run on the instances to be backed up and are responsible for all the work of actually backing themselves up.

## Configuration
When starting an agent it expects the --config flag with a path to the config file. The config is a toml document with the following:

* ip

The IP for the http api to listen on

* port

The port for the http api to listen on

* privatekey

The private key is used for signing requests sent to the Coordinator

[coordinator]

Configuration options for the cooridnator that will be managing this Agents

* address

IP:PORT for the agent: used to lookup their public key and verify the signature

* PublicKey

Public key file of the coordinator
