pipeline {
    agent any
    
    environment {
        GOPATH = "${pwd}"
    }
    
    stages {
        stage('Build') {
            steps {
                echo 'Building...'
                sh '''GOPATH = "${pwd}"
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
            }
        }
        stage('Test') {
            steps {
                sh 'Test'
            }
        }
    }
}
