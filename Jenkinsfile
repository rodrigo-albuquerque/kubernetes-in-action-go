pipeline {
    agent any
    environment {
        DOCKER_IMAGE_NAME = "digofarias/app:v1"
        registry = "digofarias/app"
        registryCredential = 'dockerhub'
        dockerImage = ''
    }
    stages {
        stage('Cloning Git Repository') {
            steps {
                git 'https://github.com/rodrigo-albuquerque/kubernetes-in-action-go.git'
            }
        }
        stage('Build Docker Image') {
            steps {
                script {
                    dockerImage = docker.build registry + ":$BUILD_NUMBER"
                }
            }
        }
        stage('Push Docker Image') {
            steps {
                script {
                    docker.withRegistry( '', registryCredential ) {
                        dockerImage.push()
                    }
                }
             }
        }
        stage('Clear Unused docker image') {
            steps {
                sh "docker rmi $registry:$BUILD_NUMBER"
            }
        }
    }
}
