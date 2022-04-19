OSTYPE ?= darwin
env ?= dev
app_name ?= hackernews-api
repo_root = ${PWD}
cluster_name ?= aws_test_cluster
namespace_name ?= hackernews-api
helm_root ?= ${PWD}/chart/hackernews-api
docker_repo ?= snigdhasambit/hackernews-api
docker_release ?= 1.0
KUBECONFIG ?= /root/.kube/config

set_namespace:
	kubectl config use-context ${cluster_name} --kubeconfig=${KUBECONFIG} \
  	kubectl config set-context ${cluster_name} --namespace ${namespace_name} --kubeconfig=${KUBECONFIG}

docker_build:
	docker build ${repo_root} -t ${docker_repo}:${docker_release}

docker_release: docker_build
	docker push ${docker_repo}:${docker_release}
docker_run:
	docker run -p 8080:8080 ${docker_repo}:${docker_release}
	# open localhost:8080/stories

deploy_dry:
	helm upgrade -i ${app_name} ${helm_root} \
    --set ImageVersion=${docker_release} \
    --debug \
    --dry-run

deploy: 
	helm upgrade -i ${app_name} ${helm_root} \
	--set ImageVersion=${docker_release} \
    --debug

service:
	minikube service hackernews-api --url

destroy:
	helm delete ${app_name}

# optional installation of prometheus

prometheus_deploy:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts \
	helm install prometheus prometheus-community/prometheus \
    kubectl expose service prometheus-server --type=NodePort --target-port=9090 --name=prometheus-server-np

prometheus_url:
	minikube service prometheus-server-np --url

destroy_prometheus:
	helm delete prometheus

