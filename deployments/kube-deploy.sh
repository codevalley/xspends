#!/bin/bash

# Exit script on any error
set -e

# Define the base name for the docker image
IMAGE_BASE_NAME="xspends-image"

# Check if a tag is provided as a command-line argument
TAG=$1

if [ -z "$TAG" ]; then
    # Generate a timestamp if no tag is provided
    TAG=$(date +%Y%m%d%H%M%S)
    echo -e "\nUse the command < ./deploy.sh custom-tag > to specify a custom tag for the build image\n"
fi

# Build and tag the Docker image with the timestamp or provided tag
echo -e "\n Set the tag $TAG..."
docker build -t ${IMAGE_BASE_NAME}:${TAG} .

# Function to update deployment yaml
update_deployment_yaml() {
    local SED_CMD=$1
    local FILE=deployments/app-deployment.yaml
    echo -e "\n Apply temp tag to deployment yaml..."
    $SED_CMD "s|${IMAGE_BASE_NAME}:TAG_PLACEHOLDER|${IMAGE_BASE_NAME}:${TAG}|g" $FILE
    echo -e "\n Delete old deployment..."
    kubectl delete -f $FILE
    echo -e "\n Apply new deployment..."
    kubectl apply -f $FILE
    echo -e "\n Remove tag to deployment yaml..."
    $SED_CMD "s|${IMAGE_BASE_NAME}:${TAG}|${IMAGE_BASE_NAME}:TAG_PLACEHOLDER|g" $FILE
}

# Detect the environment and set the sed command accordingly
case "$(uname -s)" in
    Darwin)
        SED_CMD="sed -i ''"
        ;;
    Linux|MINGW*|MSYS*|CYGWIN*)
        SED_CMD="sed -i"
        ;;
    *)
        echo "Unknown operating system"
        exit 1
        ;;
esac

update_deployment_yaml "$SED_CMD"

# Apply the Kubernetes service definition
# This is typically not needed every time unless the service definition has changed
echo -e "\n Apply the service..."
kubectl apply -f deployments/xspends-service.yaml

# Wait for a few seconds to allow the pod to start
echo "Waiting for the pod to start..."

# Optional: Add a line to confirm deployment
echo "Deployment updated with image ${IMAGE_BASE_NAME}:${TAG}"
sleep 6

# Get the Minikube service URL
minikube service xspends-service --url
