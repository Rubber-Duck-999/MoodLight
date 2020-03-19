pipeline {
    agent any
    
    environment {
        GOPATH = "${pwd}"
    }
    /*
    stages {
        stage('Build') {
            steps {
                echo 'Building...'
                sh 'cd src'
                sh 'go build'
            }
        }
        stage('Test') {
            steps {
                sh 'echo "Test"'
                sh 'cd src'
                sh 'go build'
            }
        }
    }
    post {
        failure {
            emailext body: 'Failed to build FH', subject: 'Build Failure', to: '$DEFAULT_RECIPIENTS'
        }
    }*/
}
