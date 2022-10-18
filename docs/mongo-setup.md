# Dev env setup

## Docker based Single node replica set without keys

**Step 1:** single node replica set without authentication

```sh
docker container run -d --name dev-mongo -p 27017:27017 mongo:latest --replSet devReplicaSet 
```

Note:

- By default mongod process binds to localhost, traffic is not allowed from outside world other than local

**Step 2:** Initiate replica-set and create user

```sh
docker exec -it dev-mongo bash

mongosh

# initialize replicaSet
rs.initiate()

# check replicaSet config
rs.conf()

# check replicaSet status
rs.status()

# create user
use admin
db.createUser({user:"mongoadmin", pwd:"secret", roles:["root"]})

# exit from mongosh
exit
```

**Step 3:** restart dev-mongo container

```sh
docker container restart dev-mongo 
```

**Connection string** :
To connect single node replica set, directConnection flag needed.

    mongodb://mongoadmin:secret@localhost:27017/?directConnection=true

## Three node replica set using kubernetes

**TODO**
