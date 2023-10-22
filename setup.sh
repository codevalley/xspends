#!/bin/bash

# Generate a secret key
SECRET=$(openssl rand -base64 32)

# Create a Kubernetes secret
kubectl create secret generic jwt-secret --from-literal=jwt-key=$SECRET

# Create the database and tables
mysql -h 127.0.0.1 -P 4000 -u root < setup.sql

# Deploy your application (assuming you have a deployment file)
kubectl apply -f app-deployment.yaml

