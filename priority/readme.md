## Synopsis
The `NtSim` is a `Python` implementation `of the the `Crawler` running over a random generated network. The objective is to examen the prioritization rule and the resulting crawling dynamic

## Notes
- The implementation does not yet account for network dynamics, node details and neighbors remain fixed

## Running
- Clone the repository with the existing file strucure  
- go to the top directory and from cl  run
    -   `source setup.py`, this adds the path to the library to `PYTHONPATH` 
    - Go to `test` and run `python NtSim_test.py`

The result is a list of all the nodes discovered, either connected directly or to their direct neigbor. 