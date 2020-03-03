pipeline {
    agent any
    
    environment {
        GOPATH = "${pwd}"
    }
    
    stages {
        stage('Build') {
            steps {
                echo 'Building...'
                sh './buildFaultHandler.sh'
            }
        }
        stage('Test') {
            steps {
                sh 'echo "Test"'
                sh './testFaultHandler.sh'
            }
        }
    }
    post {
        failure {
            emailext body: 'Failed to build FH', subject: 'Build Failure', to: '$DEFAULT_RECIPIENTS'
        }
    }
}
