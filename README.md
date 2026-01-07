# darkwing-go

Minimal Golang stub API and worker with CI/CD to GHCR and k3s deployment.

## Overview
- API: HTTP server on port 8080 with `/health`, `/hello`, `/version`.
- Worker: simple heartbeat logs every few seconds.
- Docker: multi-stage images published to GHCR.
- k8s: Deployments and Service (NodePort 30080) under namespace `darkwing`.
- GitHub Actions: build, push, and deploy to your k3s VDS.

## Local Dev (Windows)
Prereqs: Go 1.21+, Docker Desktop.

```powershell
go version
go run ./api
# In another terminal
go run ./worker
```

API test:
```powershell
curl http://localhost:8080/health
curl http://localhost:8080/hello
```

## Repository Setup
1. Create the public repo `petelinmn/darkwing-go` on GitHub.
2. In this folder, initialize and push:
```powershell
git init
git add .
git commit -m "Initial stub API/worker, k8s, CI"
git branch -M main
git remote add origin https://github.com/petelinmn/darkwing-go.git
git push -u origin main
```

## GHCR package visibility
Images are pushed to `ghcr.io/petelinmn/darkwing-go/{api,worker}`.
After first workflow run, set each package to Public in GitHub → Packages.

## VDS k3s access (KUBECONFIG secret)
On your VDS (2.56.127.170) with k3s installed:
1. SSH in and copy kubeconfig:
```bash
sudo cat /etc/rancher/k3s/k3s.yaml > ~/darkwing-kubeconfig.yaml
```
2. Edit the `server:` address in the file to use your public IP and port (default 6443):
```yaml
server: https://2.56.127.170:6443
```
3. Create a limited deploy user (optional, cluster-admin shown):
```bash
kubectl create serviceaccount deployer -n kube-system
kubectl create clusterrolebinding deployer-binding \
	--clusterrole cluster-admin \
	--serviceaccount kube-system:deployer
```
4. Get the token and build a kubeconfig for the deployer (optional). Alternatively, keep using `k3s.yaml`.
5. Add a GitHub secret `KUBECONFIGCONTENT` with the entire kubeconfig file contents.

## CI/CD: Build & Deploy
The workflow `.github/workflows/build-and-deploy.yml`:
- Builds `api` and `worker` images and pushes to GHCR using `GITHUB_TOKEN`.
- Uses `KUBECONFIGCONTENT` to apply manifests and roll deployments to the commit SHA.

Trigger: push to `main`.

## GitLab CI/CD (Alternative)
You can also use GitLab to build/push to its Container Registry and deploy to k3s.

1. Create a GitLab project (e.g., `gitlab.com/petelinmn/darkwing-go`) and enable Container Registry.
2. Ensure the registry allows public pulls (Project Settings → Visibility → Container Registry → Allow public pull access), or add an `imagePullSecret` in k8s.
3. Add CI variables:
	- `KUBECONFIGCONTENT`: contents of your updated kubeconfig with `server: https://2.56.127.170:6443`.
	- (Optional) If using private registry, set up k8s `imagePullSecrets` and reference them in Deployments.
4. The pipeline in `.gitlab-ci.yml` builds and pushes:
	- `registry.gitlab.com/<namespace>/<project>/api:{latest,sha-<commit>}`
	- `registry.gitlab.com/<namespace>/<project>/worker:{latest,sha-<commit>}`
5. On push to `main`, jobs:
	- Build `api` and `worker` (Docker-in-Docker) and push two tags: `latest`, `sha-<commit>`.
	- Deploy: applies manifests and sets Deployment images to `sha-<commit>`.

Manual verify:
```bash
kubectl -n darkwing get pods
curl http://2.56.127.170:30080/health
```

Notes:
- If Docker-in-Docker is unavailable on your runner, switch to Kaniko or a shell executor with Docker available.
- If you prefer GHCR from GitLab CI, set CI variables `REGISTRY=ghcr.io`, `IMAGE_PREFIX=ghcr.io/petelinmn/darkwing-go`, and adapt `.gitlab-ci.yml` to `docker login` with a GitHub PAT (write:packages).

## Kubernetes Manifests
Apply manually (optional):
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/api-deployment.yaml
kubectl apply -f k8s/api-service.yaml
kubectl apply -f k8s/worker-deployment.yaml
```

Access API:
- NodePort: http://2.56.127.170:30080/health

## Troubleshooting
- If image pull fails with `Unauthorized`, ensure GHCR packages are Public.
- If kubectl fails in CI, verify `KUBECONFIGCONTENT` secret is correct and server IP reachable.
- For Traefik/Ingress, you can add an Ingress later if you prefer hostname routing.