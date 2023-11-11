# NodeNeighborhood

## What is it?
___
Node Neighborhood is a continuous service that enumerates the global network of distributed Etherium Nodes, and maps indivudal nodes to their list of peers(neighbors).
- Each Node in the Ethereum Node network contains a list of neighbors it stores in a Distributed Hash Table.
- NodeNeighbood utilizes a list of known highly available boot nodes to obtain further lists of Neighbors, and crawls nodes in those lists, creating a list of seen nodes.
- It iterates over this list periodically for "Freshness" and updates the list as necessary, add in newly seen nodes and updating known nodes to reflect changes. 
