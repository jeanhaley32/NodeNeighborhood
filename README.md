# NodeNeighborhood

# Ethereum Network Controls

The objective of the _Node Neighborhood (NN)_ project is to implement controls for the Ethereum network, measuring critical parameters impacting the ablity of the Ethereum Blockchain to support a fair and balanced market for decentralized exchanges (DEX), based on smart contracts executing on top of the bockchain.

The controls implemented
1. Network size
2. Gossip propagation 
3. Blockchain liveness

__Network Size__: The NN system keeps track of the nodes in the Network and their neighboring relation, maintaining a database of nodes and their neighbors collected by a crawler. 

__Gossip Propagation__: monitoring the spread of transaction requests through the gossip network

__Blockchain Liveness__: The system will track the _block rate_, the rate at which new blocks are added and finalized on the blockchain. 


## Crawler Operation

Node Neighborhood is a continuous service that enumerates the global network of distributed Etherium Nodes mapping indivudal nodes to their list of peers(neighbors).

  - Each Node in the Ethereum Node network contains a list of neighbors it stores in a Distributed Hash Table.
  - NodeNeighbood utilizes a list of known highly available boot nodes to obtain further lists of Neighbors, and crawls nodes in those lists, creating a list of seen nodes.
  - It iterates over this list periodically for "Freshness" and updates the list as necessary, add in newly seen nodes and updating known nodes to reflect changes. 
