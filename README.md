# Cryptocurrency Distributed System

### Project Overview
This project is a mock cryptocurrency miners distributed system. There are three components in this system: client, pool and miner. Client is responsible for publishing task with a `data` string and an `upperBound`. Pool will receive the task from the client and divide it into multiple jobs, and then assign the jobs to different miners in the pool. Miners will find the nonce with minimum hash value within its own interval and return it in a `ProofOfWorkMsg`. Then pool returns the nonce (between 0 and `upperBound`) with the minimum hash value to the client. 

### Extra Features
- The number of intervals is based on a ratio of the current available miners and the upperBound.
- In the competition version, we used a map to store the speed of each miner. If the speed of a miner is less than the a ratio of the average speed of all miners then it will be treated as a slow miner. If a job is assigned to a slow miner, then we will reschedule this job to another miner. 

### TestCase Overview
This project contains the following tests and reaches 85% test coverage:
- client tests:
    - Test1: test whether the client address match or not
    - Test2: test whether the pool address match or not
    - Test3: test the situation when client disconnect with the pool
    - Test4: test the situation of client change
- interval tests:
    - Test1: upperBound < numIntervals (with no remainder)
    - Test2: upperBound < numIntervals ((with remainder)
    - Test3: upperBound > numIntervals
    - Test4: numIntervals < 0
- miner tests:
    - Test1: test the creation and the `Mine()` function of the miner
    - Test2: test the `Shutdown()` function of the miner
    - Test3: test the situation when a miner is shutdown
- pool tests: 
    - Test1: test the situation when several pools share a same address
    - Test2: test the situation when pool is disconnected
    - Test3: test the situation when create a pool with a invalid port number
    - Test4: test the situation when a pool connect with wrong addressed miner

### Known Bugs
Not found.