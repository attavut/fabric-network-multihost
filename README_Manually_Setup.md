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

### On VM-Org2
- sudo apt-get install unzip
- [upload crypto-config.zip & channel-artifacts.zip  via gcp console]
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

```

### On VM-Org1

##### 9). create 'ca0.example.com' by execute below command (before you do so, replace {put the name of secret key} with the name of the secret key. You can find it under '/crypto-config/peerOrganizations/org1.example.com/ca/')
```
cd ~/fabric-samples/fabric-network-multihost
```
```
sudo docker run --rm -d --network="poc-net" --name ca0.example.com -p 7054:7054 -e FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server -e FABRIC_CA_SERVER_CA_NAME=ca0.example.com -e FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca.org1.example.com-cert.pem -e FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/{put the name of secret key} -v $(pwd)/crypto-config/peerOrganizations/org1.example.com/ca/:/etc/hyperledger/fabric-ca-server-config -e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=poc-net hyperledger/fabric-ca sh -c 'fabric-ca-server start -b admin:adminpw -d'
```

##### 10). create 'orderer.example.com' by execute below command 
```
sudo docker run --rm -d --network="poc-net" --name orderer.example.com -p 7050:7050 -e ORDERER_GENERAL_LOGLEVEL=debug -e ORDERER_GENERAL_LISTENADDRESS=0.0.0.0 -e ORDERER_GENERAL_LISTENPORT=7050 -e ORDERER_GENERAL_GENESISMETHOD=file -e ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block -e ORDERER_GENERAL_LOCALMSPID=OrdererMSP -e ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp -e ORDERER_GENERAL_TLS_ENABLED=false -e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=poc-net -v $(pwd)/channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block -v $(pwd)/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp:/var/hyperledger/orderer/msp -w /opt/gopath/src/github.com/hyperledger/fabric hyperledger/fabric-orderer orderer
```

##### 11). create 'couchdb0' for peer0.org1 by execute below command 
```
sudo docker run --rm -d --network="poc-net" --name couchdb0 -p 5984:5984 -e COUCHDB_USER= -e COUCHDB_PASSWORD= -e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=poc-net hyperledger/fabric-couchdb
```

##### 12). create 'peer0.org1.example.com' by execute below command 
```
sudo docker run --rm -d --link orderer.example.com:orderer.example.com --network="poc-net" --name peer0.org1.example.com -p 7051:7051 -e CORE_LEDGER_STATE_STATEDATABASE=CouchDB -e CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb0:5984 -e CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME= -e CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD= -e CORE_PEER_ADDRESSAUTODETECT=true -e CORE_PEER_ID=peer0.org1.example.com -e CORE_PEER_ADDRESS=peer0.org1.example.com:7051 -e CORE_PEER_LISTENADDRESS=0.0.0.0:7051 -e CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org1.example.com:7051 -e CORE_PEER_LOCALMSPID=Org1MSP -e CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock -e CORE_LOGGING_LEVEL=DEBUG -e CORE_PEER_PROFILE_ENABLED=true -e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=poc-net -e CORE_PEER_TLS_ENABLED=false -e CORE_PEER_GOSSIP_USELEADERELECTION=false -e CORE_PEER_GOSSIP_ORGLEADER=true -e CORE_PEER_NETWORKID=peer0.org1.example.com -e CORE_NEXT=true -e CORE_PEER_ENDORSER_ENABLED=true -e CORE_PEER_COMMITTER_LEDGER_ORDERER=orderer.example.com:7050 -e CORE_PEER_GOSSIP_IGNORESECURITY=true -v /var/run/:/host/var/run/ -v $(pwd)/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/fabric/msp -w /opt/gopath/src/github.com/hyperledger/fabric/peer hyperledger/fabric-peer peer node start
```

Make sure all containers are up by execute below command 
```
sudo docker ps -a
```


### On VM-Org2
##### 13). create 'ca1.example.com' by execute below command (before you do so, replace {put the name of secret key} with the name of the secret key. You can find it under '/crypto-config/peerOrganizations/org2.example.com/ca/')
```
cd ~/fabric-samples/fabric-network-multihost
```
```
sudo docker run --rm -d --network="poc-net" --name ca1.example.com -p 7054:7054 -e FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server -e FABRIC_CA_SERVER_CA_NAME=ca1.example.com -e FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca.org2.example.com-cert.pem -e FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/{put the name of secret key} -v $(pwd)/crypto-config/peerOrganizations/org2.example.com/ca/:/etc/hyperledger/fabric-ca-server-config -e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=poc-net hyperledger/fabric-ca sh -c 'fabric-ca-server start -b admin:adminpw -d'
```


##### 14). create 'couchdb1' for peer0.org2 by execute below command 
```
sudo docker run --rm -d --network="poc-net" --name couchdb1 -p 6984:5984 -e COUCHDB_USER= -e COUCHDB_PASSWORD= -e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=poc-net hyperledger/fabric-couchdb
```

##### 15). create 'peer0.org2.example.com' by execute below command 
```
sudo docker run --rm -d --link orderer.example.com:orderer.example.com --link peer0.org1.example.com:peer0.org1.example.com --network="poc-net" --name peer0.org2.example.com -p 9051:9051  -e CORE_LEDGER_STATE_STATEDATABASE=CouchDB -e CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb1:5984 -e CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME= -e CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD= -e CORE_PEER_ADDRESSAUTODETECT=true -e CORE_PEER_ID=peer0.org2.example.com -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 -e CORE_PEER_LISTENADDRESS=0.0.0.0:9051 -e CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org2.example.com:9051 -e CORE_PEER_LOCALMSPID=Org2MSP -e CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock -e CORE_LOGGING_LEVEL=DEBUG -e CORE_PEER_PROFILE_ENABLED=true -e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=poc-net -e CORE_PEER_TLS_ENABLED=false -e CORE_PEER_GOSSIP_USELEADERELECTION=false -e CORE_PEER_GOSSIP_ORGLEADER=true -e CORE_PEER_NETWORKID=peer0.org2.example.com -e CORE_NEXT=true -e CORE_PEER_ENDORSER_ENABLED=true -e CORE_PEER_COMMITTER_LEDGER_ORDERER=orderer.example.com:7050 -e CORE_PEER_GOSSIP_IGNORESECURITY=true -v /var/run/:/host/var/run/ -v $(pwd)/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp:/etc/hyperledger/fabric/msp -w /opt/gopath/src/github.com/hyperledger/fabric/peer hyperledger/fabric-peer peer node start
```

Make sure all containers are up by execute below command 
```
sudo docker ps -a
```

### On VM-Org1
##### 16). up fabric network by execute below command to spawn CLI fabric-tools (This will install the CLI container and will execute the './scripts/script.sh')
```
cd ~/fabric-samples/fabric-network-multihost

sudo chmod +x ./scripts/script.sh
```

```
sudo docker run --rm -it --network="poc-net" --name cli --link orderer.example.com:orderer.example.com --link peer0.org1.example.com:peer0.org1.example.com -p 12051:7051 -p 12053:7053 -e GOPATH=/opt/gopath -e CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock -e CORE_LOGGING_LEVEL=DEBUG -e CORE_PEER_ID=cli -e CORE_PEER_ADDRESS=peer0.org1.example.com:7051 -e CORE_PEER_LOCALMSPID=Org1MSP -e CORE_PEER_TLS_ENABLED=false -e CORE_PEER_NETWORKID=cli -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp -e CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=poc-net  -v /var/run/:/host/var/run/ -v $(pwd)/../chaincode/:/opt/gopath/src/github.com/chaincode -v $(pwd)/crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ -v $(pwd)/scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/ -v $(pwd)/channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts -w /opt/gopath/src/github.com/hyperledger/fabric/peer hyperledger/fabric-tools /bin/bash -c './scripts/script.sh'
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

### For test query & invoke, you can use 'fabric client application' from  https://github.com/attavut/fabric-client-application


### For couchDB, you can view the transactions at (open it on browser)
- VM-Org1 (PC 1): http://<VM-Org1 External IP address>:5984/_utils/#/database/mychannel_poc_cc/_all_docs
- VM-Org2 (PC 2): http://<VM-Org2 External IP address>:6984/_utils/#/database/mychannel_poc_cc/_all_docs












