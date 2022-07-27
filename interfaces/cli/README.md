# Ukama CLI

## Configuration 

Every command has it's own configuration structure. That is usually located in `pkg/cmd/<command>/config.go`.
Values of config could be set via config file, environment variables of flags if `flag` tag is defined for the field.

Here is an example: 
``` go
type deployConfig struct {		
	BaseDomain       string      `flag:"baseDomain" validate:"required"`
	Helm             *HelmConfig	
}

type HelmConfig struct {
	RepoUrl string `flag:"helmRepo" default:"https://raw.githubusercontent.com/ukama/helm-charts/repo-index"`
	Token   string `flag:"token"`
}
```

We can set `BaseDomain` in three ways:
- via command line `ukama deploy --baseDomain=<domain>`.
- via environment variable `UKAMA_BASEDOMAIN=<domain>`
- using config file `.ukama.yaml` like in example below
```
deploy:
    baseDomain: example.com  
```


Here is the way to set nested value `Token`:
- via command line `ukama deploy --=<token>`.
- via environment variable `UKAMA_HELM_TOKEN=<token>`
- using config file `.ukama.yaml` like in example below.
```
deploy:
    helm:      
      token: <token>
```


# Deploy

## Deploying service 

You will need `kubectl` configured to access the cluster. 
- Create secret with access to Ukama registry in target namespace  
```
    kubectl create secret docker-registry regcred \
    --docker-server=${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com \
    --docker-username=AWS \
    --docker-password=$(aws ecr get-login-password)
```

Deploy ukama service(helm chart) by running below command
```
ukama deploy --service ukama@v0.1.149-dev --baseDomain ukama-test.com
```

## Provision Cluster
### Prerequsites 

To provision AWS cluster you will need:
- AWS CLI. Refer to [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) for installation instructions.
- Kops. Refer to [Kops](https://kops.sigs.k8s.io/getting_started/install/) for installation instructions.

Make sure you run `aws configure` to set up your credentials.


Deploying cluster:

```
ukama deploy cluster --bucket ukama-cluster-cli-denis-eu-1 --dns test-denis.k8s.local --verbose --region eu-central-1
```
