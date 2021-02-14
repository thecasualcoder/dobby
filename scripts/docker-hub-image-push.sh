#!/bin/bash

git_tag=$1

if [ -z "$git_tag" ]
then
    echo "no git tag has been provided. usage: docker-hub-image-push.sh <git-tag>"
    exit 1
fi


echo "Using the git tag '$git_tag' as dobby Docker image tag"

set -e

echo "Logging into Docker Hub using credentials..."
docker login --username thecasualcoder --password $DOCKER_HUB_TOKEN

echo "Building dobby Docker image..."
docker build -t thecasualcoder/dobby -t karuppiah7890/dobby:$git_tag .

echo "Pushing dobby Docker image..."
docker push thecasualcoder/dobby
