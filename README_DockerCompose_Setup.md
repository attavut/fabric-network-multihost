## Build Hyperledger Fabric Network (Multi-Host)

based on https://github.com/hyperledger/fabric-samples/tree/release-1.4/first-network

The directions for using the original 1 Org are documented at wahabjawed's article at Medium.
["Hyperledger Fabric on Multiple Hosts"](https://medium.com/@wahabjawed/hyperledger-fabric-on-multiple-hosts-a33b08ef24f)


#### Network Topology

So the network that we are going to build will have the following below components. For this example we are using 2 VM instances on GCP (VM-Org1 and VM-Org2):

**VM-Org1**

* A Certificate Authority (ca0) 
* An Orderer (orderer)
* 1 PEER of Org1 (peer0.org1)
* CouchDB (couchdb0)

**VM-Org2**

* A Certificate Authority (ca1) 
* 1 PEER of Org1 (peer0.org2)
* CouchDB (couchdb1)

```
  +----+-----+        +-----+----+
  |          |        |          |
+-+  VM-Org1 |        |  VM-Org2 +---+
| |          |        |          |   |
| +--------+-+        +---+------+   |
|          |              |          |
|          |              |          |
|          |              |          |
|          +-------+------+          |
|                  +                 |
|             Docker swarm           |
|           (overlay network)        |
|                                    |
+--+                                +--+
    ca0.example.com                     ca1.example.com
    orderer.example.com                 peer0.org2.example.com
    peer0.org1.example.com              couchdb1     
    couchdb0                           
       
```
## Prerequisites

##### 1). List of components that I used

Ubuntu-16.04
Fabric — 1.4.0
Docker (version 17.06.2-ce or greater)
```
Follow https://docs.docker.com/install/linux/docker-ce/ubuntu/
(
    - SET UP THE REPOSITORY (step 1-4)
    - INSTALL DOCKER CE (step 1-2)
)
```
Docker-Compose (version 1.14.0 or greater)
```
Follow https://docs.docker.com/compose/install/
(
    - Install Compose >> Linux (step 1-2)
)
```

#### 2). Make sure specified ports are not blocked with firewall. 
###### Ports for docker swarm
- tcp:2377, tcp:7946
- udp:7946, udp:4789

###### Ports for Fabric containers
- ca0.example.com, ca1.example.com >>> port 7054
- orderer.example.com >>> port 7050    
- peer0.org1.example.com >>> port 7051
- peer0.org2.example.com >>> port 9051
         



## Setup Process
### On all hosts (VM-Org1, VM-Org2)
##### 1). Install Samples project, Binaries and Docker Images)
```
sudo curl -sSL http://bit.ly/2ysbOFE | bash -s 1.4.0
```
##### 2). go to fabric-samples directory
```
cd fabric-samples
```
##### 3). clone Fabric Network MulitHost source code from  https://github.com/attavut/fabric-network-multihost
```
sudo git clone https://github.com/attavut/fabric-network-multihost.git
```

### On VM-Org1
##### 4). Initialize a swarm
```
sudo docker swarm init --advertise-addr <VM-Org1 External IP address>
```
##### 5). Join the swarm with the other host as a manager (VM-Org1 will create swarm)
```
sudo docker swarm join-token manager
```
It will output something like this
```
docker swarm join --token SWMTKN-1-xxxxx8kgzlalp0d3udtaz2jaavvp5d4xg7tyr0g5vhfm8pwpm5-8ckx8yq0r5a3dyyyyy <VM-Org1 External IP address>:2377
```

### On VM-Org2
##### 6) VM-Org2 join swarm
We will copy output from 5). (the one on your terminal, not the one above) and add '--advertise-addr' for VM-Org2. Execute it on VM-Org2 terminal to make it join swarm
```
docker swarm join --token SWMTKN-1-xxxxx8kgzlalp0d3udtaz2jaavvp5d4xg7tyr0g5vhfm8pwpm5-8ckx8yq0r5a3dyyyyy <VM-Org1 External IP address>:2377 --advertise-addr <VM-Org2 External IP address>
```


### On VM-Org1
##### 7). Create docker network (overlay network) to attach all Hyperledger services ("poc-net" in my case).

```
sudo docker network create --attachable --driver overlay poc-net
```
##### 8). Generate Fabric Crypto Material ("crypto-config" and "channel-artifacts") 
```
cd ~/fabric-samples/fabric-network-multihost/

sudo chmod +x ./byfn.sh

sudo ./byfn.sh generate

```
This will generate Crypto Material for you in "crypto-config" and "channel-artifacts" folder. You must copy these folders to VM-Org2 (copy to all hosts as you want)

Example copy step & command
```
### On VM-Org1
- sudo apt-get install zip
- cd ~/fabric-samples/fabric-network-multihost/crypto-config
- sudo zip -r crypto-config.zip .
- [download via gcp console]

- cd ~/fabric-samples/fabric-network-multihost/channel-artifacts
- sudo zip -r channel-artifacts.zip .
- [download via gcp console]

- cd ~/fabric-samples/fabric-network-multihost/
- [download 'docker-compose-org2.yaml' via gcp console]

### On VM-Org2
- sudo apt-get install unzip
- [upload 'crypto-config.zip' & 'channel-artifacts.zip' & 'docker-compose-org2.yaml'  via gcp console]
- sudo rm -rf ~/fabric-samples/fabric-network-multihost/crypto-config
- sudo mkdir ~/fabric-samples/fabric-network-multihost/crypto-config
- sudo cp ~/crypto-config.zip ~/fabric-samples/fabric-network-multihost/crypto-config
- cd ~/fabric-samples/fabric-network-multihost/crypto-config
- sudo unzip crypto-config.zip -d .

- sudo rm -rf ~/fabric-samples/fabric-network-multihost/channel-artifacts
- sudo mkdir ~/fabric-samples/fabric-network-multihost/channel-artifacts
- sudo cp ~/channel-artifacts.zip ~/fabric-samples/fabric-network-multihost/channel-artifacts
- cd ~/fabric-samples/fabric-network-multihost/channel-artifacts
- sudo unzip channel-artifacts.zip -d .

- sudo rm -rf ~/fabric-samples/fabric-network-multihost/docker-compose-org2.yaml
- sudo cp ~/docker-compose-org2.yaml ~/fabric-samples/fabric-network-multihost

```

### On VM-Org1
##### 9). create all org1's containers by execute below command
```
cd ~/fabric-samples/fabric-network-multihost
```
```
sudo ./byfn.sh org1up
```


### On VM-Org2
##### 10). create all org2's containers by execute below command
```
cd ~/fabric-samples/fabric-network-multihost

sudo chmod +x ./byfn.sh
```
```
sudo ./byfn.sh org2up
```

### On VM-Org1
##### 16). up fabric network by execute below command to spawn CLI fabric-tools (This will install the CLI container and will execute the './scripts/script.sh')
```
cd ~/fabric-samples/fabric-network-multihost

sudo chmod +x ./scripts/script.sh
```

```
sudo ./byfn.sh up
```

If you see this, it means that the script has been executed

```
.
.
.

"========= All GOOD, BYFN execution completed =========== "


" _____   _   _   ____   "
"| ____| | \ | | |  _ \  "
"|  _|   |  \| | | | | | "
"| |___  | |\  | | |_| | "
"|_____| |_| \_| |____/  "

```

This script ('./scripts/script.sh') will:
- Create channel (mychannel in this case)
- Make peer0.org1 and peer0.org2 join the channel.
- UpdateAnchorPeers for peer0.org1 and peer0.org2
- InstallChaincode for peer0.org1 and peer0.org2
- InstantiateChaincode for peer0.org1
- Query test on peer0.org1
- Invoke test
- Query test on peer0.org2















