# Jean Crawler proposal. 

The Crawler I have in mind should have a series of features.

   Each Node will be it's own ENR record, and a neighbors list. 
   ```
   type Node struct {
    record       enr       // base64 encoding RLP representation 'enr:...'
    neighbors    []enr     // List of Node's neighbors (DHT)
    LastSeen     time.Time // Time Node was last seen
    Freshness    uint      // Contrived freshness number used to determine if the node should be discarded. 
   }
   ```

1. Maintains an externally accessible database of Etherium Nodes and their Neighbors. 
2. Starting with a list of BootNodes, creates an initial set of known Etherium adresses.
3. Crawls that list of nodes, generating a record that contains it's [ENR](https://eips.ethereum.org/EIPS/eip-778) and it's neighbors(It's portion of the [DHT](https://docs.ipfs.tech/concepts/dht)).

4. Discerns a Level of "Freshness" for each node, removing them from the Database when there is indication that the node
   is no longer available. 
5. if a node hasn't been seen in X length of time, searches for that node again, updating if any information has changed and altering it's freshness accordingly. 
6. This crawler will run continuously, creating a solid representation of whole Etherium node network. 

## Initializer
The Initializer is the Boot Straping mechanism for the crawler.

It is responsible for the following:

1. Starts,and maintains in a WorkGroup the Work Delegator, Freshness Tracker, and Node Ingestor. 
2. Creates the Channels needed for each operation to communicate with each other.
3. Converts initial BootNodes into work, and passes them to the Work Delegator.
3. manages over-all routines, gracefully exiting program as routines shutdown, or restarts them if necessary.

## Work Delegator
The Work Delegator monitors the Work Delegator channel (pending name change)
It is responsible for the following: 

1. Obtains work, sees what action is being requested, and spawns a new worker to execute that action. 

The worker is very simple and straightforward. It is designed to generate laborers to perform functions based
on objects passed to it that contain the requesite action. 

## Worker
A worker is an ephemeral Go Routine spawned from the Work Delegator to perform a task. 

- Spawned from a work item, a worker is provided with an "action" to perform, and the appropriate contextual variables
  Needed to succesfully execute that function. 
- Upon execution, the result is bundles and sent down a destination channel. 
- Upon a Receipt request, the worker shuts down. 

 > For the purpose of this crawler, the action is going to be to start up a ephemeral Etherium Node, and obtain the nodes most updated information, and as many neighbors as it can associate with that node. It then communicates this to the Node Ingestor to act upon that information. 

 By designing this in a non-specific way, we can utilize workers in the future to perform multiple task types, leaving us open to implement new features. 


## Node Ingestor
 The Node Ingestor listens for new nodes passing to it on the Node Ingestion Channel. 
 
 It is responsible for the following:
 - Fields New Nodes from the Node Ingestor Channel.
    0. if the Worker reports that a Node was unable to be contacted, it will decrement freshness of a node. 
    1. Checks to see if node has been seen previously and if so, if the sequence number has changed.  
        1. If node has been seen, and the sequence number is the same as what is listed in "Seen Nodes".
            1. It will increment freshness for that node in the "seen list" by 1. 
            (in the case that the node's sequence number is in fact lower, we should probably not increment freshness.)
        2. If the Node has been seen, and the sequence number is higher. 
            1. It will update the Seen Nodes sequence number, increment freshness, and send a write request to the database.
    2. If the Node has not been seen before, it will create a new entry in seen nodes, and will add a write request to the Database path
        To add the new entry. 
    3. For each of the Neighbors in the nodes list of neighbors, it will check against "seen nodes", and pass on as new work, any network that isn't seen. 

- Handles Requests from the Freshness Tracker.
    1. Provides a list of seen nodes
    2. Decrements freshness of specific "seen nodes" upon request
    3. Sends request to DB for node removal, and pops node from list of seen nodes. 

## Freshness tracker. 
Every node should have a "freshness" measurement. 

The purpose of freshness is to quantify a nodes availability/reliability.

> Each Node won’t necessarily respond to a request every time one is made, this doesn’t mean it no longer exists. At the same time, in order to keep a reliable and up-to-date view of the Ethereum Network, it’s necessary to set a threshold for node responsiveness. This threshold helps in determining the minimum number of responses a node should provide within a certain timeframe to be considered active in the network.

*I'm not sure what the best method of execution for this concept is. My initial idea is to create an integer that increases by one every time a node is seen, and decrements by a percentage everytime we fail to communicate with that node.*

The Freshness tracker will be responsible for this, and will also be responsible for creating new work if we get into a state where none is being generated from the node ingester. 

1. Every X window of time (5, 10 minutes, or even less) the tracker will request a full copy of the known nodes list. 
2. we will check the last time each node in the list was seen, and based on the delta between then and now, we derive how many runs went by where the node was not seen.
3. Freshness is then decremented by a percentage derived from that length of time.
> an example distribution could be 1:0% 2:5% 3:20% 4:40% 5: 80% 6: 100%. We round up to ensure at least 1 point is removed every time. This distribution of percentages i arbitrary, and is not based on any insightful data. 

4. If freshness drops below a target value, or if it is reduced to zero The Freshness Tracker will send a removal request back to the ingestor. 
5. If it doesn't fall to zero, or below that target threshold, then we send a request to the ingestor to udpate the freshness value to reflect. 
4. For Every Node that hasn't been checked since the last epoch (Time in which freshness tracker runs), we pass a new work item to the work delegator. This way we continuosly check the Ethereum network for changes.


### Extra Credit
---
1. I'd like to make a metrics path that derives values from work items. I.E. average work time, how much work performed every second etc. This would involve obtaining information from work items, and crunching those numbers elswhere, maybe logging them. We can use this information to create a metrics dashboard that can be used to monitor the health of the service. 

2. Right now, this is going to be a monolithic applicaiton, comprised of multiple goroutines. In the future, it may be fun to set this up as a series of microservices communicating through GRPC protobufs. I'm not sure how that would change things, or if there is even a benefit to this decentralization of labor, but it would be an educationation experience. 
