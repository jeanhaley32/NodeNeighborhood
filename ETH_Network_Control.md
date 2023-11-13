# Ethereum Network Controls

The objective of the _Node Neighborhood (NN)_ project is to implement controls for the Ethereum network, measuring critical parameters impacting the ablity of the Ethereum Blockchain to support a fair and balanced market for decentralized exchanges (DEX), based on smart contracts executing on top of the bockchain.

The controls implemented
1. Network size
2. Gossip propagation 
3. Blockchain liveness

__Network Size__: The NN system keeps track of the nodes in the Network and their neighboring relation, maintaining a database of nodes and their neighbors collected by a crawler. 

__Gossip Propagation__: monitoring the spread of transaction requests through the gossip network

__Blockchain Liveness__: The system will track the _block rate_, the rate at which new blocks are added and finalized on the blockchain. 

__