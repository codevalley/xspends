General instructions to setup the entire setup in a local minikube.
Assuming you have installed
- docker
- helm
- minikube.


$minikube start --cpus=6 --memory='12288mb' 

$eval $(minikube docker-env) 

$helm repo add pingcap https://charts.pingcap.org/ 
$helm update
$kubectl create namespace tidb-cluster
$helm install tidb-operator pingcap/tidb-operator --version v1.5.1 --namespace tidb-cluster
$kubectl create -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.5.1/manifests/crd.yaml

$helm install tidb-cluster pingcap/tidb-cluster --version v1.5.1 -f values-tidb.yaml --namespace tidb-cluster

$kubectl port-forward svc/tidb-cluster-tidb 4000:4000 -n tidb-cluster

$mysql -h 127.0.0.1 -P 4000 -u root

$docker build -t xpends-image .

$./setup.sh

$kubectl apply -f app-deployment.yaml

$kubectl apply -f xpends-service.yaml

$minikube service xpends-service --url

Verify the service by checking register and login curls

$curl -X POST -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpass"}' <MINIKUBE_SERVICE_URL>/register
Replace <MINIKUBE_SERVICE_URL> with the URL you got in the previous step.

$curl -X POST -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpass"}' <MINIKUBE_SERVICE_URL>/login


some useful kubectl commands
$kubectl get pods
$kubectl describe pod <POD_NAME>
$kubectl logs <POD_NAME> -c xpends-container
