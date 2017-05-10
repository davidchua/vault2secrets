# Vault2Secrets

![vault2secrets architecture](https://cloud.githubusercontent.com/assets/68039/25882064/5712396c-3573-11e7-8e62-7c9a8c46290b.png)

### Objectives

Vault2Secrets is developed to solve the problem of syncing secrets stored on Hashicorp's Vault with Kubernetes.

### Problem

Operators want to allow Developers to self-deploy their apps onto Kubernetes without them having access to sensitive data that they need to use in the normal functioning of their apps.

Operators would store such sensitive data into an external resource like Hashicorp Vault and like them to be automatically injected as Environment Variables within those deployed application pods.

### Function

Vault2Secrets is a Kubernetes Controller that can be deployed on a Kubernetes Cluster which makes use of a ThirdPartyResources called `CustomSecret` to retrieve Hashicorp Vault stored data and convert them into secure Kubernetes Secret Objects.

Operators can then reference these Secret Objects to be loaded onto their deployment's environment variables.

Whenever the Vault Secret that is being monitored has a change, Vault2Secrets will be able to automatically pick it up and update the respective Vault Secret Objects.


## Instructions

Important Environment Flags
```
NAMESPACE=default
KUBERNETES_API_ENDPOINT=127.0.0.1:8001
```

### Developing

To get a local copy of the controller running, please run

`make`

This will pull all the necessary dependencies and build a binary in the current directory

### Deploying

To deploy, you can find the example scripts in `examples/`

TL;DR

1. `kubectl create -f examples/tpr.yml`
2. `kubectl create -f examples/vault2secrets.yml`
3. Modify `examples/generic-secret.yml` with your `VAULT TOKEN` and deploy
4. Modify `examples/example-custom-secret.yml` and deploy

### Credits

This controller wouldn't be possible if not for @kelseyhightower [Kubernetes Certificate Manager](https://github.com/kelseyhightower/kube-cert-manager) and his presentation at [PuppetConf 2016](https://www.youtube.com/watch?v=HlAXp0-M6SY)
