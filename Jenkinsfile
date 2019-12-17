pipeline {
    agent any
    
    environment {
        GOPATH = "${pwd}"
    }
    
    stages {
        stage('Build') {
            steps {
                echo 'Building...'
                sh 'cd src'
                sh 'go version'
                sh 'go get -v github.com/streadway/amqp'
                sh 'go get -v github.com/sirupsen/logrus'
                sh 'go get -v github.com/scorredoira/email'
                sh 'go get -v gopkg.in/yaml.v2'
                sh 'go get -v github.com/akamensky/argparse'
                sh 'go get -v github.com/clarketm/json'
                sh 'ls'
                sh 'pwd'
                sh 'go install github.com/Rubber-Duck-999/exeFaultHandler'
                sh 'go get -u -v github.com/golang/lint/golint'
            }
        }
        stage('Test') {
            steps {
                sh 'Test'
            }
        }
    }
}
