#!/bin/bash

# Generate a secret key
SECRET=$(openssl rand -base64 32)

# Create a Kubernetes secret
kubectl create secret generic jwt-secret --from-literal=jwt-key=$SECRET
kubectl create secret generic db-credentials --from-literal=DB_DSN="root:@tcp(tidb-cluster-tidb.tidb-cluster.svc.cluster.local:4000)/xspends"

# Create the database and tables
mysql -h 127.0.0.1 -P 4000 -u root < ./scripts/setup.sql

# Deploy your application (assuming you have a deployment file)
kubectl apply -f deployments/app-deployment.yaml

