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
    echo "\nUse the command < ./deploy.sh custom-tag > to specify a custom tag for the build image\n"
fi

# Build and tag the Docker image with the timestamp or provided tag
echo "\n Set the tag $TAG..."
docker build -t ${IMAGE_BASE_NAME}:${TAG} .

# Update the Kubernetes deployment file with the new image tag
# The following sed command works for macOS. If you are on Linux, you may need to remove the '' after -i
echo "\n Apply temp tag to deployment yaml..."
sed -i '' "s|${IMAGE_BASE_NAME}:TAG_PLACEHOLDER|${IMAGE_BASE_NAME}:${TAG}|g" deployments/app-deployment.yaml

# Delete the old Kubernetes deployment
echo "\n Delete old deployment..."
kubectl delete -f deployments/app-deployment.yaml

# Apply the updated Kubernetes deployment
echo "\n Apply new deployment..."
kubectl apply -f deployments/app-deployment.yaml

echo "\n Remove tag to deployment yaml..."
sed -i '' "s|${IMAGE_BASE_NAME}:${TAG}|${IMAGE_BASE_NAME}:TAG_PLACEHOLDER|g" deployments/app-deployment.yaml

# Apply the Kubernetes service definition
# This is typically not needed every time unless the service definition has changed
echo "\n Apply the service..."
kubectl apply -f deployments/xspends-service.yaml

# Wait for a few seconds to allow the pod to start
echo "Waiting for the pod to start..."

# Optional: Add a line to confirm deployment
echo "Deployment updated with image ${IMAGE_BASE_NAME}:${TAG}"
sleep 6

# Get the Minikube service URL
minikube service xspends-service --url

