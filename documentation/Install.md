# Installation
You can easily install LogBook via Go CLI:

    go get github.com/alexgunkel/logbook
    cd ${GOPATH}/src/github.com/alexgunkel/logbook
    go install

## Requirements
This program requires Golang at least in version 1.6 (the least version
that's tested by me. For the tests to run you need at least version 1.7.

## Usage
Start LogBook by typing

    logbook

*LogBook* by default listens on port 8080. To change this just set the
environmental variable PORT to any other value. Likewise you can configure
the hostname:

    PORT=80 HOST=127.0.0.1 logbook

will change the port to 80 and the hostname to 127.0.0.1.

## Docker and Orchestration
There will be a Docker image available which you can get with

    docker pull alexandergunkel/logbook

For further information visit the
[public repository](https://hub.docker.com/r/alexandergunkel/logbook/)
at [Docker Hub](https://hub.docker.com)