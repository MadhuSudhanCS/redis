redis
========

* In-Memory Data Store.
* Implements subset of Redis commands.
* Command reply is according to Redis Protocol.

Commands supported
------------------

* GET
* GETBIT
* QUIT
* SAVE
* SET
* SETBIT
* ZADD
* ZCARD
* ZCOUNT
* ZRANGE


Description
-----------

* Redis server listens to TCP port 15000 for all the incomming requests.
* Server loads all the key-value pairs to cache from the db file at the start up.
* During exit(CTRL+C/SIGINT) cache is saved to json file.

Cache
-----

* Cache is made up of buckets.
* Buckets hold key-value pair.
* To distribute the cache load, number of buckets in the cache can be increased/decreased which is conifigurable.

Prerequisites
-------------

* Install GO 1.5+ - https://golang.org/dl/
* Add GO to the path

Installation
------------

$ go get 'github.com/madhusudhancs/redis'
$ cd $GOPATH/src/github.com/madhusudhancs/redis
* Build the server 
    * go build
* To list the options available
    * ./redis -h
* To start the server 
    * ./redis -f <Absolute_PATH_TO_JSON_FILE>
        * ./redis -f ./db/dump.json
        * Ensure dir db is created in the path

