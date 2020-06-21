#!/usr/bin/env groovy

pipeline {
    agent {
    kubernetes {
      yaml """
apiVersion: v1
kind: Pod
metadata:
  labels:
    some-label: some-label-value
spec:
  containers:
  - name: docker
    image: docker:18.05-dind
    securityContext:
        privileged: true
  - name: golang
    image: golang
    command: ["sleep", "10000"]
"""
        }
    }

    environment {
        registry = "dineshba/endgame"
        registryCredential = 'dockerhub'
        dockerImage = ""
    }

    stages {
        stage('Validations') {
            steps {
                container('golang') {
                    script {
                        sh "make test"
                    }
                }
            }
        }   

        stage('docker build') {
            steps {
              container('docker') {
                script {
                    def imageTag = "$registry:latest"
                    dockerImage = docker.build imageTag
                }
              }
            }
        }

        stage('docker publish') {
            steps {
              container('docker') {
                script {
                    docker.withRegistry("", registryCredential) {
                        dockerImage.push()
                    }
                }
              }
            }
        }

        stage('Deploy App') {
             steps {
                script {
                    kubernetesDeploy(configs: "dobby.yaml", kubeconfigId: "mykubeconfig")
                }
            }
        }
    }
}