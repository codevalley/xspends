
# XSpends Project Setup Guide

This guide provides step-by-step instructions to set up the XSpends project in a local Minikube environment.

## Prerequisites
Ensure you have the following tools installed:
- Docker
- Helm
- Minikube

## Setup Instructions

1. **Start Minikube**:
   ```bash
   minikube start --cpus=6 --memory='12288mb'
   ```

2. **Configure Docker to use Minikube's Docker daemon**:
   ```bash
   eval $(minikube docker-env)
   ```

3. **Add the PingCAP Helm repository and update it**:
   ```bash
   helm repo add pingcap https://charts.pingcap.org/
   helm repo update
   ```

4. **Set up the TiDB Cluster**:
   ```bash
   kubectl create namespace tidb-cluster
   helm install tidb-operator pingcap/tidb-operator --version v1.5.1 --namespace tidb-cluster
   kubectl create -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.5.1/manifests/crd.yaml
   helm install tidb-cluster pingcap/tidb-cluster --version v1.5.1 -f deployments\values-tidb.yaml --namespace tidb-cluster
   ```
   **Make sure the clusters are running (it will take a few minutes)**:
   
   ```bash
   kubectl get pods --namespace tidb-cluster -l app.kubernetes.io/instance=tidb-cluster
   The output will look something like this in the beginning
   ```bash
   tidb-cluster-pd-0                         0/1     ContainerCreating
   tidb-cluster-pd-1                         0/1     ContainerCreating
   ...
   ```
   Should change to something like this...
   ```bash
   tidb-cluster-pd-0                         1/1     Running   0          27m
   tidb-cluster-pd-1                         1/1     Running   0          27m
   tidb-cluster-pd-2                         1/1     Running   0          27m
   tidb-cluster-tidb-0                       2/2     Running   0          25m
   tidb-cluster-tidb-1                       2/2     Running   0          25m
   tidb-cluster-tikv-0                       1/1     Running   0          26m
   tidb-cluster-tikv-1                       1/1     Running   0          26m
   tidb-cluster-tikv-2                       1/1     Running   0          26m
   ```
   **Check DB cluster status and troubleshoot if needed**:
   ```bash
   kubectl get pods --namespace tidb-cluster -l app.kubernetes.io/instance=tidb-cluster
   kubectl describe pod <pod-name> --namespace tidb-cluster
   kubectl logs <pod-name> --namespace tidb-cluster
   kubectl get pvc --namespace tidb-cluster
   kubectl get pv

   ```

5. **Port-forward TiDB service for local access**:
   ```bash
   kubectl port-forward svc/tidb-cluster-tidb 4000:4000 -n tidb-cluster
   ```

6. **Access the TiDB database**:
   ```bash
   mysql -h 127.0.0.1 -P 4000 -u root
   ```

7. **Build the application's Docker image**:
   ```bash
   docker build -t xspends-image .
   ```

8. **Run the setup script**:
   ```bash
   ./scripts/setup.sh
   ```

9. **Deploy the application**:
   ```bash
   kubectl apply -f deployments/app-deployment.yaml
   kubectl apply -f deployments/xspends-service.yaml
   ```

10. **Access the application's service URL**:
    ```bash
    minikube service xspends-service --url
    ```

11. **Verify the service**:
    Replace `<MINIKUBE_SERVICE_URL>` with the URL obtained in the previous step:
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpass"}' <MINIKUBE_SERVICE_URL>/register
    curl -X POST -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpass"}' <MINIKUBE_SERVICE_URL>/login
    ```

## Useful `kubectl` Commands

- List all pods: `kubectl get pods`
- Describe a specific pod: `kubectl describe pod <POD_NAME>`
- View logs for a specific container: `kubectl logs <POD_NAME> -c xspends-container`

## Commands After Code Changes

If you make code changes and wish to redeploy:

1. Rebuild the Docker image:
   ```bash
   docker build -t xspends-image .
   ```

2. Redeploy the application:
   ```bash
   kubectl delete -f app-deployment.yaml
   kubectl apply -f app-deployment.yaml
   ```

3. Verify pods are running:
   ```bash
   kubectl get pods
   ```

4. Access the updated application's service URL:
   ```bash
   minikube service xspends-service --url
   ```

5. Quickly verify an endpoint (after replacing `<MINIKUBE_SERVICE_URL>`):
   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpass"}' <MINIKUBE_SERVICE_URL>/login
   ```
   or 

   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpass"}' <MINIKUBE_SERVICE_URL>/register
   ```



6. Check logs for any issues:
   ```bash
   kubectl get pods # to get the pod name
   kubectl logs <POD_NAME> -c xspends-container
   ```
7. Docker userful commands
   ```
   docker images | grep xspends-image
   # Optionally, clean up old Docker images
   docker image prune -a -f
   ```
### Continuous builds and deploys
There is a simple script which can be used to rebuild and deploy images for continuous testing. It is avalable in `deployments` folder. 
You can run the following command from project root folder. 
If you are booting up the docker fresh following these commands
```bash
$minikube start --cpus=6 --memory='12288mb'
$eval $(minikube docker-env)
$kubectl port-forward svc/tidb-cluster-tidb 4000:4000 -n tidb-cluster #this would consume the terminal. 
$./deployments/kube-deploy.sh 
#if you need to flush the DB data
$mysql -h 127.0.0.1 -P 4000 -u root < ./scripts/flush_data.sql
```
The script is self explanatory, you can review it to understand how it works. 

### Testing and mocks

1. Generate mocks
```
$ mockgen -source=./kvstore/rawkv_interface.go -destination=./kvstore/mock/mock_rawkv.go -package=mock
$ mockgen -source=./models/db.go -destination=./models/mock/mock_DBExecutor.go -package=mock
```