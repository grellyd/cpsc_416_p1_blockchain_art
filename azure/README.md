# AZURE VMs

We have 7 virtual machines up on Azure: 1 server, 1 miner with no art nodes, 3 miners with one art node, and 2 miners with three art nodes. The miners use the naming scheme Miner<#ArtNodes>Art<#>

The *username* for each VM is `supfoo`, and they live within the `supfoo416` resource group. The server names and their IP addresses are as follows:

Server - 40.65.104.57
Miner0Art1 - 40.65.107.115
Miner1Art1 - 40.65.107.130
Miner1Art2 - 40.65.124.64
Miner1Art3 - 40.65.107.162
Miner3Art1 - 40.65.104.179
Miner3Art2 - 40.65.108.120

Server runs on port 12345,
Miners run on port 8000,
and Art nodes run on a randomly generated port.

VMs shut off at 1am PST daily.

So long as you've had your public SSH key added onto the machine, you should be able to ssh into a machine using the command `ssh -i ~/.ssh/id_rsa supfoo@<IP-ADDRESS>`.

# Scripts
*NOTE*: Scripts need to be run from _INSIDE_ the azure folder!

`azureinstall.sh` - bash script given to us for setting up a Go environment

`start_server.sh` - used to start the server

`start_miner.sh` - used to start a miner. Usage is `start_miner.sh <NAME_OF_MACHINE>` For example, `./start_miner.sh Miner1Art2`. If you are running the server locally, add a second argument: `./start_miner.sh Miner1Art2 local`.

# Running an art node/test client
There is currently no script for running an art node.

However, it's pretty simple:

1) Run a server
2) Run a miner
3) Run the following from within the top-level directory:

go run art-app.go <MINER_IP>:8000 `cat azure/keyfiles/MINER_NAME.priv` `cat azure/keyfiles/MINER_NAME.pub`

