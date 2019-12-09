pipeline {
    agent any
    environment {
        DOCKER_IMAGE_NAME = "digofarias/app:v1"
        registry = "digofarias/app"
        registryCredential = 'dockerhub'
        dockerImage = ''
        PROJECT_ID = 'rodrigo-albuquerque'
        CLUSTER_NAME = 'rodrigo-k8s-cluster'
        LOCATION = 'europe-west1-d'
        CREDENTIALS_ID = 'kubernetes-demo'
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
        stage('Deploy to GCP Kubernetes cluster') {
            steps {
                step([
                $class: 'KubernetesEngineBuilder',
                projectId: env.PROJECT_ID,
                clusterName: env.CLUSTER_NAME,
                location: env.LOCATION,
                manifestPattern: 'k8s-dc-v1.yaml',
                credentialsId: env.CREDENTIALS_ID,
                verifyDeployments: true])
                }
            }
      }
}
