#!/bin/bash

# Function to create a cluster
create_cluster() {
  local cluster_name=$1
  local config=$2
  
  echo "Creating cluster: $cluster_name..."
  cat <<EOF > ${config}
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: ${cluster_name}
nodes:
  - role: control-plane
EOF

  if [ "$cluster_name" == "cluster-one" ]; then
    echo "  - role: worker" >> ${config}
    echo "  - role: worker" >> ${config}
  elif [ "$cluster_name" == "cluster-two" ]; then
    echo "  - role: worker" >> ${config}
    echo "  - role: worker" >> ${config}
    echo "  - role: worker" >> ${config}
    echo "  - role: worker" >> ${config}
  fi

  kind create cluster --config ${config}
  echo "Cluster $cluster_name created successfully."
}

# Create cluster-one
create_cluster "cluster-one" "kind-cluster-one.yaml"

# Create cluster-two
create_cluster "cluster-two" "kind-cluster-two.yaml"

# Verify clusters
echo "Verifying clusters..."
kind get clusters

# Deploy app to a cluster
deploy_app() {
  local cluster_context=$1
  echo "Deploying app to $cluster_context..."
  kubectl config use-context ${cluster_context}
  kubectl create ns development
  kubectl config set-context --current --namespace=development
  kubectl apply -f project/chart/configmap.yaml
  kubectl apply -f project/chart/deployment.yaml
  kubectl apply -f project/chart/service.yaml
  echo "App deployed to $cluster_context successfully."
}

# Deploy app to cluster-one
deploy_app "kind-cluster-one"

# Deploy app to cluster-two
deploy_app "kind-cluster-two"

# Create load balancer cluster
echo "Creating load balancer cluster..."
create_cluster "load-balancer" "kind-load-balancer.yaml"

# Deploy load balancer
echo "Deploying load balancer..."
kubectl config use-context kind-load-balancer
kubectl create ns development
kubectl config set-context --current --namespace=development
kubectl apply -f clusterLB/nginx-configmap.yaml

echo "Validating load balancer setup..."
kubectl get services
kubectl port-forward service/nginx-api-gateway 9092:80 &
echo "Load balancer deployed and port-forwarding started."

# Set up Prometheus for cluster-one
echo "Setting up Prometheus for cluster-one..."
kubectl config use-context kind-cluster-one
kubectl create -f prometheus-setup/ns.yaml
kubectl config set-context --current --namespace=monitoring

# Edit external labels to match the cluster name
echo "Applying Prometheus configuration for cluster-one..."
kubectl apply -f prometheus-setup/configmap.yaml
kubectl apply -f prometheus-setup/deployment.yaml
kubectl apply -f prometheus-setup/service.yaml
kubectl apply -f prometheus-setup/rbac.yaml
kubectl apply -f prometheus-setup/thanus-sidecar-svc.yaml

# Repeat Prometheus setup for cluster-two
echo "Setting up Prometheus for cluster-two..."
kubectl config use-context kind-cluster-two
kubectl create -f prometheus-setup/ns.yaml
kubectl config set-context --current --namespace=monitoring

echo "Applying Prometheus configuration for cluster-two..."
kubectl apply -f prometheus-setup/configmap.yaml
kubectl apply -f prometheus-setup/deployment.yaml
kubectl apply -f prometheus-setup/service.yaml
kubectl apply -f prometheus-setup/rbac.yaml
kubectl apply -f prometheus-setup/thanus-sidecar-svc.yaml

# Create and deploy Thanos cluster
echo "Creating Thanos cluster..."
create_cluster "thanos-cluster" "kind-thanos-cluster.yaml"
kubectl config use-context kind-thanos-cluster
kubectl create -f thanus-setup/ns.yaml
kubectl config set-context --current --namespace=monitoring

echo "Validating Thanos setup..."
kubectl get services
kubectl port-forward service/thanos-querier 9090:9090 &
echo "Thanos setup complete and port-forwarding started."

echo "All steps completed successfully!"
