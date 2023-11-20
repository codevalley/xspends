#!/bin/bash

# Generate a secret key
SECRET=$(openssl rand -base64 32)

# Create a Kubernetes secret
kubectl create secret generic jwt-secret --from-literal=jwt-key=$SECRET
kubectl create secret generic db-credentials --from-literal=DB_DSN="root:@tcp(tidb-cluster-tidb.tidb-cluster.svc.cluster.local:4000)/xspends?parseTime=true"
#for production set these values. For local (minikube) testing we dynamically get the value.
kubectl create secret generic app-host --from-literal=app-host="127.0.0.1"
kubectl create secret generic app-port --from-literal=app-port="8080"

#connection params
kubectl create secret generic DB_MAX_OPEN_CONNS --from-literal=DB_MAX_OPEN_CONNS="25"
kubectl create secret generic DB_MAX_IDLE_CONNS --from-literal=DB_MAX_IDLE_CONNS="25"
kubectl create secret generic DB_CONN_MAX_LIFETIME --from-literal=DB_CONN_MAX_LIFETIME="5"
# Create the database and tables
mysql -h 127.0.0.1 -P 4000 -u root < ./scripts/setup.sql
#for windows
#Get-Content .\scripts\setup.sql | mysql -h 127.0.0.1 -P 4000 -u root

# Deploy your application (assuming you have a deployment file)
kubectl apply -f deployments/app-deployment.yaml

