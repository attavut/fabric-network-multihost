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


## Setup Process

### Manually Setup
[a relative link](README_Manually_Setup.md)

### Using Docker Compose 
[a relative link](README_DockerCompose_Setup.md)














