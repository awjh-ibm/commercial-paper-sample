# commercial-paper-sample

This sample covers a simple commercial paper contract using the updated fabric programming model.

## Using this sample in chaincode developer mode
Copy this folder into the chaincode folder of fabric-samples/chaincode-docker-devmode.

Make a clone of fabric (which includes the changes made in CR 25801) in your gopath at $GOPATH/src/github.com/hyperledger/fabric and run `make docker`.

Update the docker-compose-simple.yaml file in fabric-samples/chaincode-docker-devmode to:

```
version: '2'

services:
  orderer:
    container_name: orderer
    image: hyperledger/fabric-orderer:amd64-latest
    environment:
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=orderer
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=orderer.block
      - ORDERER_GENERAL_LOCALMSPID=DEFAULT
      - ORDERER_GENERAL_LOCALMSPDIR=/etc/hyperledger/msp
      - GRPC_TRACE=all=true,
      - GRPC_VERBOSITY=debug
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./msp:/etc/hyperledger/msp
      - ./orderer.block:/etc/hyperledger/fabric/orderer.block
    ports:
      - 7050:7050
  peer:
    container_name: peer
    image: hyperledger/fabric-peer:amd64-latest
    environment:
      - CORE_PEER_ID=peer
      - CORE_PEER_ADDRESS=peer:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer:7051
      - CORE_PEER_LOCALMSPID=DEFAULT
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=DEBUG
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp
    volumes:
        - /var/run/:/host/var/run/
        - ./msp:/etc/hyperledger/msp
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start --peer-chaincodedev=true -o orderer:7050
    ports:
      - 7051:7051
      - 7053:7053
    depends_on:
      - orderer

  cli:
    container_name: cli
    image: hyperledger/fabric-tools:amd64-latest
    tty: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=DEBUG
      - CORE_PEER_ID=cli
      - CORE_PEER_ADDRESS=peer:7051
      - CORE_PEER_LOCALMSPID=DEFAULT
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp
    working_dir: /opt/gopath/src/chaincodedev
    command: /bin/bash -c './script.sh'
    volumes:
        - /var/run/:/host/var/run/
        - ./msp:/etc/hyperledger/msp
        - ./chaincode:/opt/gopath/src/chaincodedev/chaincode
        - ./:/opt/gopath/src/chaincodedev/
        - $GOPATH/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric
    depends_on:
      - orderer
      - peer

  chaincode:
    container_name: chaincode
    image: hyperledger/fabric-ccenv
    tty: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=DEBUG
      - CORE_PEER_ID=example02
      - CORE_PEER_ADDRESS=peer:7051
      - CORE_PEER_LOCALMSPID=DEFAULT
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp
    working_dir: /opt/gopath/src/chaincode
    command: /bin/bash -c 'sleep 6000000'
    volumes:
        - /var/run/:/host/var/run/
        - ./msp:/etc/hyperledger/msp
        - ./chaincode:/opt/gopath/src/chaincode
        - $GOPATH/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric
    depends_on:
      - orderer
      - peer
```

Startup the docker containers:
`docker-compose -f docker-compose-simple.yaml up`

Connect to the Chaincode container:
`docker exec -it chaincode bash`

Inside the docker container build the go and start the chaincode:
```
cd commercial_paper
go build
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./commercial_paper
```

Connect to the CLI container:
`docker exec -it cli bash`

Install the chaincode:
`peer chaincode install -p chaincodedev/chaincode/commercial_paper -n mycc -v 0`

Instantiate the chaincode:
`peer chaincode instantiate -n mycc -c '{"Args":["contract_Setup"]}' -C myc -v 0`

Create a paper:
`peer chaincode invoke -n mycc -c '{"Args":["contract_CreatePaper", "PAPER1", "20", "1000"]}' -C myc`

Create a second paper:
`peer chaincode invoke -n mycc -c '{"Args":["contract_CreatePaper", "PAPER2", "20", "1000"]}' -C myc`

Add the two papers to the market "US_BLUE_ONE" with the discount "20":
`peer chaincode invoke -n mycc -c '{"Args":["contract_ListOnMarket", "US_BLUE_ONE", "20", "PAPER1", "PAPER2"]}' -C myc`

Check that they have been added to the market:
`peer chaincode query -n mycc -c '{"Args":["contract_RetrieveMarket", "US_BLUE_ONE"]}' -C myc`
