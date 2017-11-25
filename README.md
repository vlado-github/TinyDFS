TinyDFS
=======

A learning project.
Building a simple distributed database from ground up.

- Database is key-value based and uses files for storage.
- Database can also be used as general file distributed storage. 
- The file data should be kept in more than one node. 
- Implements a consistency policy for reads and writes. 
- Tackless the problem of how keeping data in more than one node helps to provide fault tolerance. 

## Plans

This project needs to provide evaluation results in terms of performance improvement due to storing data as well as the improved fault tolerance due to multi-node system.

## Build

Follow these steps:
- clone or extract downloaded archive
- run terminal
- run command
```bash
go get github.com/google/uuid
```
- position to **tinydfs** directory
- run command
```bash
go build
```

## Run

Follow these steps:
- position to **tinydfs** directory
- run command
```bash
./tinydfs master
```
- open another terminal
- run command (note: omit 'master' argument)
```bash
./tinydfs
```

## Tests

Follow these steps:
- position to **messaging** directory
- or position to **persistance** directory
- and run command:
```bash
go test -v
```

## Ideas for ToDo

- implement interface for write/read for different DBs support
- research and introduce vector clocks support
- add consistency hashing for both reads and writes
- support for configurable number of replicas and partitions
- research and introduce protocol for master selection
- Benchmarking: writes/reads for n-nodes (LAN, web)

## Literature
- Google File System. Sanjay Ghemawat, Howard Gobioff, Shun-Tak Leung. SOSP 2003. Student Presenter: Rita Chiu.
- Dynamo: Amazon’s Highly Available Key-value Store. Giuseppe DeCandia, Deniz Hastorun, Madan Jampani, Gunavardhan Kakulapati, Avinash Lakshman, Alex Pilchin, Swaminathan Sivasubramanian, Peter Vosshall and Werner Vogels. SOSP 2007.

