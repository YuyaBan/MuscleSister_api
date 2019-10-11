# MuscleSister_api
## Description
* A Simple Http Api Server made in Go (Golang).
* Now, not linked DB Server.
    * TODO: link DB Server.

## How to build in localhost
* You can deploy easily via Dockerfile.
1. docker build -t musclesister_api . 
2. docker run -e "PORT=<port_number>" -p <port_number>:<port_number> -t musclesister_api

## Usage example (2019/10/11)
* Get
    * curl localhost:<port_number>/1
* POST
    * curl -XPOST localhost:<port_number>/1 -d '{"id":1,"name":"foo","user":"bar","done":false}'
* PUT
    * curl -XPUT localhost:<port_number>/1
        * PutMethod change done_flag to true
* DELETE
    * curl -XDELETE localhost:<port_number>/1