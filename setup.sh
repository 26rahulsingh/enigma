#!/bin/bash

echo "Setting up the network.."
echo "Creating channel genesis block.."

# Create the channel
docker exec -e "CORE_PEER_LOCALMSPID=RTAMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/rta.enigma.com/users/Admin@rta.enigma.com/msp" -e "CORE_PEER_ADDRESS=peer0.rta.enigma.com:7051" cli peer channel create -o orderer.enigma.com:7050 -c enigmachannel -f /etc/hyperledger/configtx/enigmachannel.tx

sleep 5

echo "EnigmaChannel genesis block created."

echo "peer0.rta.enigma.com joining the channel..."

# Join peer0.rta.enigma.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=RTAMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/rta.enigma.com/users/Admin@rta.enigma.com/msp" -e "CORE_PEER_ADDRESS=peer0.rta.enigma.com:7051" cli peer channel join -b enigmachannel.block

echo "peer0.rta.enigma.com joined the channel"

echo "peer0.insurancecompany.enigma.com joining the channel..."

# Join peer0.insurancecompany.enigma.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=InsuranceCompanyMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/insurancecompany.enigma.com/users/Admin@insurancecompany.enigma.com/msp" -e "CORE_PEER_ADDRESS=peer0.insurancecompany.enigma.com:7051" cli peer channel join -b enigmachannel.block

echo "peer0.insurancecompany.enigma.com joined the channel"

echo "peer0.resalebroker.enigma.com joining the channel..."

# Join peer0.resalebroker.enigma.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=ResaleBrokerMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/resalebroker.enigma.com/users/Admin@resalebroker.enigma.com/msp" -e "CORE_PEER_ADDRESS=peer0.resalebroker.enigma.com:7051" cli peer channel join -b enigmachannel.block

sleep 5

echo "peer0.resalebroker.enigma.com joined the channel"

echo "Installing enigma chaincode to peer0.rta.enigma.com..."

# install chaincode
# Install code on rta peer
docker exec -e "CORE_PEER_LOCALMSPID=RTAMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/rta.enigma.com/users/Admin@rta.enigma.com/msp" -e "CORE_PEER_ADDRESS=peer0.rta.enigma.com:7051" cli peer chaincode install -n enigmacc -v 1.0 -p github.com/enigma/go -l golang

echo "Installed enigma chaincode to peer0.rta.enigma.com"

echo "Installing enigma chaincode to peer0.insurancecompany.enigma.com...."

# Install code on insurancecompany peer
docker exec -e "CORE_PEER_LOCALMSPID=InsuranceCompanyMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/insurancecompany.enigma.com/users/Admin@insurancecompany.enigma.com/msp" -e "CORE_PEER_ADDRESS=peer0.insurancecompany.enigma.com:7051" cli peer chaincode install -n enigmacc -v 1.0 -p github.com/enigma/go -l golang

echo "Installed enigma chaincode to peer0.insurancecompany.enigma.com"

echo "Installing enigma chaincode to peer0.resalebroker.enigma.com..."

# Install code on resalebroker peer
docker exec -e "CORE_PEER_LOCALMSPID=ResaleBrokerMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/resalebroker.enigma.com/users/Admin@resalebroker.enigma.com/msp" -e "CORE_PEER_ADDRESS=peer0.resalebroker.enigma.com:7051" cli peer chaincode install -n enigmacc -v 1.0 -p github.com/enigma/go -l golang

sleep 5

echo "Installed enigma chaincode to peer0.resalebroker.enigma.com"

# Instantiating enigma chaincode...
echo "Instantiating enigma chaincode.."

docker exec -e "CORE_PEER_LOCALMSPID=RTAMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/rta.enigma.com/users/Admin@rta.enigma.com/msp" -e "CORE_PEER_ADDRESS=peer0.rta.enigma.com:7051" cli peer chaincode instantiate -o orderer.enigma.com:7050 -C enigmachannel -n enigmacc -l golang -v 1.0 -c '{"Args":[""]}' -P "OR ('RTAMSP.member','InsuranceCompanyMSP.member','ResaleBrokerMSP.member')"

echo "Instantiated enigma chaincode."

echo "Following is the docker network....."

docker ps
