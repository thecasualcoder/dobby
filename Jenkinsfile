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
                    def imageTag = "dineshba:${GIT_BRANCH}-${BUILD_NUMBER}"
                    sh "docker build -f Dockerfile . -t ${imageTag}"
                }
              }
            }
        }
    }
}