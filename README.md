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
5. Add a GitHub secret `KUBECONFIG_CONTENT` with the entire kubeconfig file contents.

## CI/CD: Build & Deploy
The workflow `.github/workflows/build-and-deploy.yml`:
- Builds `api` and `worker` images and pushes to GHCR using `GITHUB_TOKEN`.
- Uses `KUBECONFIG_CONTENT` to apply manifests and roll deployments to the commit SHA.

Trigger: push to `main`.

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
- If kubectl fails in CI, verify `KUBECONFIG_CONTENT` secret is correct and server IP reachable.
- For Traefik/Ingress, you can add an Ingress later if you prefer hostname routing.