# Vault2Secrets

### Objectives

Vault2Secrets is developed to solve the problem of syncing secrets stored on Hashicorp's Vault with Kubernetes.

Its objective is to take secrets stored in Hashicorp Vault and return it as Kubernetes Secret Objects.

### Problem

Operators want to allow Developers to self-deploy their apps onto Kubernetes without them having access to sensitive data that they need to use in the normal functioning of their apps.

Operators would store such sensitive data into an external resource like Hashicorp Vault and like them to be automatically injected as Environment Variables within those deployed application pods.

### Function

Vault2Secrets is a Kubernetes Controller that can be deployed on a Kubernetes Cluster which makes use of Custom ThirdPartyResources like `CustomSecret` to instruct the retrieval and access of Hashicorp Vault stored data and convert them into secure Kubernetes Secret Objects.

Operators can then instruct as part of the application's deployment manifest file to pull environment variables as `SecretRef`.

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

To deploy, run the following:

`make deploy`

This will create a Docker image where the binary will be built and stored.

You can then take the image and push it to a Docker Registry of choice.

