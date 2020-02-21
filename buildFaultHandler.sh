#!/bin/sh


export GOPATH=$PWD
echo $GOPATH
cd src

go get -v github.com/streadway/amqp
go get -v github.com/sirupsen/logrus
go get -v github.com/scorredoira/email
go get -v gopkg.in/yaml.v2
go get -v github.com/akamensky/argparse
go get -v github.com/clarketm/json
go get -v github.com/sfreiberg/gotwilio
go install github.com/Rubber-Duck-999/exeFaultHandler
#go test -cover github.com/Rubber-Duck-999/...
go test -coverprofile coverage.html github.com/Rubber-Duck-999/...
