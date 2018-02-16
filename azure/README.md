# AZURE VMs

We have 7 virtual machines up on Azure: 1 server, 1 miner with no art nodes, 3 miners with one art node, and 2 miners with three art nodes. The miners use the naming scheme Miner<#ArtNodes>Art<#>

The *username* for each VM is `supfoo`, and they live within the `supfoo416` resource group. The server names and their IP addresses are as follows:

Server - 40.65.122.229
Miner0Art1 - 20.190.43.114
Miner1Art1 - 20.190.41.8
Miner1Art2 - 20.190.43.178
Miner1Art3 - 40.65.113.111
Miner3Art1 - 40.65.117.234
Miner3Art2 - 40.65.102.54

VMs shut off at 1am PST daily.

So long as you've had your public SSH key added onto the machine, you should be able to ssh into a machine using the command `ssh -i ~/.ssh/id_rsa supfoo@<IP-ADDRESS>`.

# Setup scripts

`azureinstall.sh` - bash script given to us for setting up a Go environment
