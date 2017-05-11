#Getting Started

These are the short steps to get your Vault2Secrets running on your Kubernetes Cluster

## Install CustomSecret Third Party Resource

```
$ kubectl create -f examples/tpr.yml
```

## Deploy Vault2Secrets

```
$ kubectl create -f examples/vault2secrets.yml
```

> **Note**: `Metadata.Labels.Namespace` should be filled up. This contains the information where the controller will be installed and where it will also look for the CustomSecrets that you'll be deploying

> **Important**: At this stage, you should already have your vault running and the pod that you'll be injecting Vault to be running.

## Store your VAULT_TOKEN as a Secret

```
apiVersion: v1
kind: Secret
metadata:
  name: vault-secret
  namespace: default
type: Opaque
data:
  token: "VAULT_TOKEN"
```

You will need to create a secret with the key `token` that contains the VAULT_TOKEN which will be used as reference in the next step.

> **Note**: Please make sure they are in the same namespace as everything else.

## Launch a CustomSecret

`CustomSecrets` contains the specification that Vault2Secrets need to pull and deploy the secrets from Vault.

A typical CustomSecret will look like:

```
apiVersion: "cubiclerebels.com/v1"
kind: CustomSecret
metadata:
  name: hello
  namespace: default
spec:
  url: "http://localhost:8200"
  path: "secret/hello"
  tokenRef: "vault-secret" # name of the Kubernetes Secret where it contains the token to auth with VAULT
  secret: "name-of-secret" # name of the Kubernetes Secret in which the returned data will be stored
```

```
# spec definition
url - contains the URL to your Hashicorp Vault instance
path - secret path to pull the data from.
tokenRef - contains the name of the Kubernetes Secret which it will pull VAULT_TOKEN from
secret - the name of the output Kubernetes Secret in which the pulled data will be stored
```

You can find an example of a `CustomSecret` at `examples/example-custom-secret.yml`

## Finally

If all goes well, you should see a new `Secret` called `name-of-secret` in your defined namespace. If you do not see a `Secret` created within 1 minute, do check the `vault2secrets` logs.
