# TinyDFS

A learning project.
TinyDFS is a distributed log storage.

- Database is key-value and uses files for storage.
- The file data are kept in more than one node.
- Consistency over all nodes is eventual.
- Distributed system has built-in fail over in case of a master queue fails.

## Plans

- All messages in a system must be encrypted.
- Different databases on persistence layer could be supported if necessary interfaces are provided.
- This project needs to provide evaluation results in terms of performance improvement due to storing data as well as the improved fault tolerance due to multi-node system.

## Build

Run build command within the root repository directory (use Linux bash/Git bash):

```bash
make
```

## Run

Run command within the root repository directory:

```bash
./build/tinydfs -listen <port> -connect <ip_address> <port>
```

## Tests

Run command within the root repository directory:

```bash
make test
```

## Ideas

- implement interface for write/read for different DBs/S3 support
- research and introduce vector clocks support
- add consistency hashing for both reads and writes
- support for configurable number of replicas and partitions
- research and introduce protocol for master selection
- Benchmarking: writes/reads for n-nodes (LAN, web)

## Literature

- Google File System. Sanjay Ghemawat, Howard Gobioff, Shun-Tak Leung. SOSP 2003. Student Presenter: Rita Chiu.
- Dynamo: Amazon’s Highly Available Key-value Store. Giuseppe DeCandia, Deniz Hastorun, Madan Jampani, Gunavardhan Kakulapati, Avinash Lakshman, Alex Pilchin, Swaminathan Sivasubramanian, Peter Vosshall and Werner Vogels. SOSP 2007.
