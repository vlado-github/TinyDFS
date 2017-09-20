TinyDFS
=======

A learning project.
Building a simple distributed file system from ground up.

- The file system should have the equivalent of the open, close, read and write calls. 
- The file data should be cooperatively cached and kept in more than one node. 
- Implements a consistency policy for reads and writes. 
- Tackless the problem of how keeping data in more than one node helps to provide fault tolerance. 

## Plans

This project needs to provide evaluation results in terms of performance improvement due to caching as well as the improved fault tolerance due to multi-node caching.

## Run

Follow these steps:
- clone or extract downloaded archive
- run terminal
- run: go get github.com/google/uuid
- position to TinyDFS/node_app
- run "go build" command
- run "./node_app master"

- run another terminal (same steps)
- but run "./node_app" instead (omit 'master' argument)

## Ideas for ToDo

- handling connections closing
- implement write persistance (file)
- research and introduce vector clocks support
- add consistency hashing for both reads and writes
- support for configurable number of replicas and partitions
- research and introduce Gossip protocol for master selection
- Benshmarking: writes/reads for n-nodes

## Literature
- Google File System. Sanjay Ghemawat, Howard Gobioff, Shun-Tak Leung. SOSP 2003. Student Presenter: Rita Chiu.
- Dynamo: Amazon’s Highly Available Key-value Store. Giuseppe DeCandia, Deniz Hastorun, Madan Jampani, Gunavardhan Kakulapati, Avinash Lakshman, Alex Pilchin, Swaminathan Sivasubramanian, Peter Vosshall and Werner Vogels. SOSP 2007.

