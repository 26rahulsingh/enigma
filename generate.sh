#!/bin/bash

rm -R crypto-config/*

./bin/cryptogen generate --config=crypto-config.yaml

rm config/*

./bin/configtxgen -profile EnigmaOrgOrdererGenesis -outputBlock ./config/genesis.block

./bin/configtxgen -profile EnigmaOrgChannel -outputCreateChannelTx ./config/enigmachannel.tx -channelID enigmachannel
