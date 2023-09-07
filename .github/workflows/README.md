# CI process for the repo 

## Convention 
- Every service should be placed in a separate directory in the root of this repo
- Directory name will be used as a service name during deployment so keep it lowercase, short and try to use lowercase letters and dashes only 
- Every service directory should contain makefile and Dockerfile
- Makefile should contain `build` and `test` steps

This convention will help us keep the CI/CD process unified, automated and maintainable. 

## Usage 

To add a Github Action's workflow for a new service all you need is run `gen-workflows.sh` script that uses `workflow-template.yaml.templ` file as a template. 

If you need a change in a workflow, do it in template and regenerate all workflows using `gen-workflows` script

In order to enable automated deployments you will need to:
1. Create and Docker repository with the service name (aka directory name)
2. Add service's chart to an [Ukama or UkamaX umbrella charts](https://github.com/ukama/helm-charts) 
3. Add entry to a helmfile in [infra-as-code repo](https://github.com/ukama/infra-as-code). Check readme there for more info


