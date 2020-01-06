# HouseGuard-FaultHandler

HouseGuard-FaultHandler is the Golang components that is part of the HouseGuard solution. 
It's role to respond to faults on the system and report them to the operator.

## Installation

This can be installed on multiple OS

```bash
sh '''export GOPATH="${PWD}"
    cd src
    go version
    go get -v github.com/streadway/amqp
    go get -v github.com/sirupsen/logrus
    go get -v github.com/scorredoira/email
    go get -v gopkg.in/yaml.v2
    go get -v github.com/akamensky/argparse
    go get -v github.com/clarketm/json
    pwd
    go install github.com/Rubber-Duck-999/exeFaultHandler
    go get -u -v github.com/golang/lint/golint
'''
```


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://github.com/Rubber-Duck-999/HouseGuard-FaultHandler/blob/master/LICENSE.txt)
