# Audit Server

The Audit server is the part of the system that will save us when things go wrong.

It logs everything so that when things stop working (database failed mid-transaction, server going down etc...)
it will be able to sift through the logs and figure out what was happening at the time of failure

## Building and using the Docker image

```
cd setup_audit_server
./setup_image.sh
```

This will create an image called audit_server. To use it:

```
docker run audit_server
```

## Building/Using the Cassandra instance

[Go here](/setup_audit_cassandra)

## Downloading the Image

```
docker pull asinha94/seng468_cassandra_audit
```