# Awesome Lab

A lab environment for testing and learning Kubernetes concepts by creating clusters, deploying applications, and setting up monitoring tools.

---

## 1. Setting Up Clusters

### Creating Two Clusters for the Application

#### **Cluster One**: 1 master node and 2 worker nodes

Use the `kind-cluster.yaml` configuration to create Cluster One:
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: cluster-one
nodes:
  - role: control-plane
  - role: worker
  - role: worker
```
Command:
```bash
kind create cluster --config kind-cluster.yaml
```

#### **Cluster Two**: 1 master node and 4 worker nodes

Update `kind-cluster.yaml` for Cluster Two:
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: cluster-two
nodes:
  - role: control-plane
  - role: worker
  - role: worker
  - role: worker
  - role: worker
```
Command:
```bash
kind create cluster --config kind-cluster.yaml
```

#### Verify Clusters:
```bash
kind get clusters
# Output:
# cluster-one
# cluster-two
```

---

## 2. Deploying the Golang App on Both Clusters

### Cluster One Deployment:
```bash
kubectl config use-context kind-cluster-one
kubectl create ns development
kubectl config set-context --current --namespace=development
kubectl apply -f project/chart/configmap.yaml
kubectl apply -f project/chart/deployment.yaml
kubectl apply -f project/chart/service.yaml
```

### Cluster Two Deployment:
```bash
kubectl config use-context kind-cluster-two
kubectl create ns development
kubectl config set-context --current --namespace=development
kubectl apply -f project/chart/configmap.yaml
kubectl apply -f project/chart/deployment.yaml
kubectl apply -f project/chart/service.yaml
```

---

## 3. Creating a Load Balancer to Share Load Between Clusters

### Create Load Balancer Cluster with 2 Nodes

Use the following configuration:
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: load-balancer
nodes:
  - role: control-plane
  - role: worker
  - role: worker
```
Command:
```bash
kind create cluster --config kind-load-balancer.yaml
```

### Deploy the Load Balancer:
```bash
kubectl config use-context kind-load-balancer
kubectl create ns development
kubectl config set-context --current --namespace=development
kubectl apply -f clusterLB/nginx-configmap.yaml
```

### Validate:
```bash
kubectl get services
# Output:
# NAME                TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
# nginx-api-gateway   NodePort   10.96.185.140   <none>        80:32632/TCP   5m27s

kubectl port-forward service/nginx-api-gateway 9092:80
# Forwarding from 127.0.0.1:9092 -> 80
# Forwarding from [::1]:9092 -> 80
```

Access the app via the load balancer:
```bash
curl --location 'localhost:9092/mirror?message=HelloWorld'
# Response:
# {
#     "mirrored_message": "HelloWorld"
# }
```

---

## 4. Setting Up Prometheus for Each Cluster

### Cluster One:
```bash
kubectl config use-context kind-cluster-one
kubectl create -f prometheus-setup/ns.yaml
kubectl config set-context --current --namespace=monitoring

# Edit external labels to match the cluster name:
# external_labels:
#   cluster: "cluster-one" 
#   region: "us-east"

kubectl create -f prometheus-setup/configmap.yaml
kubectl create -f prometheus-setup/deployment.yaml
kubectl create -f prometheus-setup/service.yaml
kubectl create -f prometheus-setup/rbac.yaml
kubectl create -f prometheus-setup/thanus-sidecar-svc.yaml
```

Repeat the same steps for other clusters.

---

## 5. Setting Up Thanos to Aggregate Metrics from Both Clusters

### Create Thanos Cluster:

Use the following configuration:
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: thanos-cluster
nodes:
  - role: control-plane
```
Command:
```bash
kind create cluster --config kind-cluster.yaml
```

### Deploy Thanos:
```bash
kubectl config use-context kind-thanos-cluster
kubectl create -f thanus-setup/ns.yaml
kubectl config set-context --current --namespace=monitoring
```

### Validate:
```bash
kubectl get services
# Output:
# NAME             TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
# thanos-querier   NodePort   10.96.253.67   <none>        9090:32743/TCP   5m50s

kubectl port-forward service/thanos-querier 9090:9090
```

Access Thanos via:
```
http://localhost:9090/targets
