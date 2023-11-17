# (Jean H) Node Neighborhood crawler Design Proposal

The Crawler I have in mind should have a series of features.

Each Node will be its own ENR record, and a neighbors list, plus information useful for making decisions based on its.

```
type Node struct {
 record       enr       // base64 encoding RLP representation 'enr:...'
 neighbors    []enr     // List of Node's neighbors (DHT)
 LastSeen     time.Time // Time Node was last seen
 Freshness    uint      // Contrived freshness number used to determine if the node should be discarded.
}

```

1. Maintains an externally accessible database of Ethereum Nodes and their Neighbors.
2. Starting with a list of Boot Nodes, creates an initial set of known Ethereum addresses.
3. Crawls that list of nodes, generating a record that contains it's [ENR](https://eips.ethereum.org/EIPS/eip-778) and its neighbors (It's portion of the [DHT](https://docs.ipfs.tech/concepts/dht)).
4. Discerns a Level of "Freshness" for each node, removing them from the Database when there is indication that the node.
is no longer available.
5. if a node hasn't been seen in X length of time, searches for that node again, updating if any information has changed and altering its freshness accordingly.
6. This crawler will run continuously, creating a solid representation of whole Ethereum node network.

## Initializer

The Initializer is the Boot Strapping mechanism for the crawler.

It is responsible for the following:

1. Starts, and maintains in a Work Group the Work Delegator, Freshness Tracker, and Node Ingester.
2. Creates the Channels needed for each operation to communicate with each other.
3. Converts initial Boot Nodes into work and passes them to the Work Delegator.
4. manages over-all routines, gracefully exiting program as routines shutdown, or restarts them if necessary.

## Work Delegator

The Work Delegator monitors the Work Delegator channel (pending name change)
It is responsible for the following:

1. Obtains work, sees what action is being requested, and spawns a new worker to execute that action.

The worker is very simple and straightforward. It is designed to generate laborers to perform functions based.
on objects passed to it that contain the requisite action.

## Worker

A worker is an ephemeral Go Routine spawned from the Work Delegator to perform a task.

- Spawned from a work item, a worker is provided with an "action" to perform, and the appropriate contextual variables.
Needed to successfully execute that function.
- Upon execution, the result is bundled and sent down a destination channel.
- Upon a Receipt request, the worker shuts down.

> For the purpose of this crawler, the action is going to be to start up a ephemeral Ethereum Node, and obtain the nodes most updated information, and as many neighbors as it can associate with that node. It then communicates this to the Node Ingester to act upon that information.
> 

By designing this in a non-specific way, we can utilize workers in the future to perform multiple task types, leaving us open to implement new features.

## Work

The Concept of work is an idea worth its own segment. 

the struct design for work would look something like this:

```go
type work struct {
   uid      uint // a unique identifier for this piece of work.
   state    bool  // 0 is not started 1 is finished
   start    time.Time // Time Work Started
   finish   time.Time // Time Work Finished
   chbundle type channelBundle struct { // Bundle of channels for communication.
        workComm ch // Channel for communication with the work routine
        target   ch // target channel for payload to be sent to after work is completed.
        parent   ch // channel to send result.
         }
   payload  type payload struct { // payload is related to the action, and result.
      context []interface // context is a set of variable needed to perform the 
                        // action. I need to think about how this is going to work. 
      action  func(context)success bool, result []byte {}
      result  []byte
    success  bool // 0 is a failure 1 is a success
   	}
}
```

A work item bundles information regarding an item of work, and the operation that work is to perform. It gets spawned off into a worker. The worker takes in the ‘payload’ and executes the action within the payload, utilizing variable loaded into “context”. It returns the result as raw []byte, to be ingested by the target path. 

Work items are comprised of

- UID
    - A unique identifier for this work, having a unique way to identify individual instances of work may be necessary for certain functions, like killing hung processes, anything that may require a direct query.
    - The “worker” process will also adopt this identifier. This way you can query your list of workers and obtain it’s comm channel.
- start
    - Start Time, so we can derive statics and see what time this piece of work was started.
    - There may be value in differentiating “work generation time” and “work execution start time.”
- finish
    - Time that the worker finished its operation.
    - There may also be value in adding a time to indicate when work was completed and successfully ingested.
- channel bundle - A bundle of channels used for communication.
    - worker channel: // potentially not necessary
        - This would be the channel used by the worker for direct communication.
        - I don’t have an operation that requires direct communication with a worker at the moment, this may change in the future, and probably should.
    - target channel: → points to a target ingester for the resulting work.
        - The target channel for the work item to be sent to for appropriate ingestion once that work is completed.  This channel being set within the item of work allows us to ingest different types of work to different paths. This would ideally be set by the delegator who would own a list of work types and their requisite paths.
    - parent channel: // potentially not necessary
        - Communication channel that points to the delegator process that spawned the work item. This is only useful if the delegator asks for a response to indicate the job itself has started. This may be necessary for communicating work health. But also, may not be necessary.
- payload - Payload represents the actual work to be performed, and the result of that work.
    - Context
        - This concept needs to be fleshed out. The idea I have is that context is going to bundle a series of variables to build the necessary contextual environment for the action to be performed by the worker successfully. Each action type will have a list of necessary “context” and the processes before the work item gets to the delegator will apply that context. The delegator can drop the work item if it deems that not enough context is applied to perform the action.
    - Action
        - This is a function that takes in “context” and returns a bool, and a payload. What actions this function performs are defined separately. This way a worker can take on whatever intention the applied action gives.
    - Result
        - This is the resulting payload from the action. This payload will be handled by the target ingester.  This is going to be a raw set of []bytes, the worker itself should have no car for what format the data is in, or what it translates to. It’s up to the ingester to care about how this data is processed once it receives it.
    - Success
        - A bool to indicate whether the requested action returns a success or a failure.

## Node Ingester

The Node ingester listens for new nodes passing to it on the Node Ingestion Channel.

It is responsible for the following:

- Fields New Nodes from the Node ingester Channel.
    1. Checks to see if node has been seen previously and if so, if the sequence number has changed.
        1. If node has been seen, and the sequence number is the same as what is listed in "Seen Nodes".
            1. It will increment freshness for that node in the "seen list" by 1.
            (In the case that the node's sequence number is in fact lower, we should probably not increment freshness.)
        2. If the Node has been seen, and the sequence number is higher.
            1. It will update the Seen Nodes sequence number, increment freshness, and send a write request to the database.
    2. If the Node has not been seen before, it will create a new entry in seen nodes, and will add a write request to the Database path to add the new entry.
    3. For each of the Neighbors in the nodes list of neighbors, it will check against "seen nodes", and pass on as new work, any network that isn't seen.
- Handles Requests from the Freshness Tracker.
    1. Provides a list of seen nodes.
    2. Decrements freshness of specific "seen nodes" upon request
    3. Sends request to DB for node removal, and pops node from list of seen nodes.

## Freshness tracker.

Every node should have a "freshness" measurement.

The purpose of freshness is to quantify a nodes availability/reliability.

> Each Node won’t necessarily respond to a request every time one is made, this doesn’t mean it no longer exists. At the same time, in order to keep a reliable and up-to-date view of the Ethereum Network, it’s necessary to set a threshold for node responsiveness. This threshold helps in determining the minimum number of responses a node should provide within a certain timeframe to be considered active in the network.
> 

*I'm not sure what the best method of execution for this concept is. My initial idea is to create an integer that increases by one every time a node is seen, and decrements by a percentage every time we fail to communicate with that node.*

The Freshness tracker will be responsible for this and will also be responsible for creating new work if we get into a state where none is being generated from the node ingester.

1. Every X window of time (5, 10 minutes, or even less) the tracker will request a full copy of the known nodes list.
2. we will check the last time each node in the list was seen, and based on the delta between then and now, we derive how many runs went by where the node was not seen.
3. Freshness is then decremented by a percentage derived from that length of time.

> an example distribution could be 1:0% 2:5% 3:20% 4:40% 5: 80% 6: 100%. We round up to ensure at least 1 point is removed every time. This distribution of percentages i arbitrary and is not based on any insightful data.
> 
1. If freshness drops below a target value, or if it is reduced to zero The Freshness Tracker will send a removal request back to the ingester.
2. If it doesn't fall to zero, or below that target threshold, then we send a request to the ingester to update the freshness value to reflect.
3. For Every Node that hasn't been checked since the last epoch (Time in which freshness tracker runs), we pass a new work item to the work delegator. This way we continuously check the Ethereum network for changes.

### Extra Credit

---

1. I'd like to make a metrics path that derives values from work items. I.E. average work time, how much work performed every second etc. This would involve obtaining information from work items, and crunching those numbers elsewhere, maybe logging them. We can use this information to create a metrics dashboard that can be used to monitor the health of the service.
2. Right now, this is going to be a monolithic application comprised of multiple goroutines. In the future, it may be beneficial to set this up as a series of microservices communicating through [GRPC ProtoBuffs](https://grpc.io/docs/what-is-grpc/introduction/). I'm not sure how that would change things, or if there is even a benefit to this decentralization of labor, but it would be an educational experience.

### Caveats

---

There is a gap in the design philosophy between which go routines run as a single operation, and which divide its work amongst workers.

The assumption I'm running off of is that the only items of work that will actually create viable bottlenecks are the spawning of ephemeral nodes for the discovery of other nodes. once obtained, everything else should be fairly quick, as work will be done through a decision tree and all operations require at most the modification of a locally stored table. I have made the database operations vague on purpose.

The idea would be to pass responsibility down to a whole new path that will then take on the burden of modifying the database and exposing that database to the outside world.

Depending on how this crawler will communicate with the database, this operation may also be time intensive. Leaving the worker's job as non-specific  `we pass an action, and contextual variables` we could make communication as easy as bundling write requets into work items, and having the work delegator delegate this job to a worker. I'm unsure if this is a reasonable idea. we may run into problem with multiple workers attempting to write similar changes to the database at once. This may not be an issue if we can adequately hash each change and compare it to already stored values to see if there is any change. Or it may be enough to compare sequence numbers.
