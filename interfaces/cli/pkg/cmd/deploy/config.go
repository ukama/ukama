package deploy

import "github.com/ukama/ukama/interfaces/cli/pkg"

type deployConfig struct {
	pkg.GlobalConfig `mapstructure:",squash"`
	Cloud            string `flag:"cloud"`
	Aws              *AwsConfig
	Service          []string    `flag:"service"`
	BaseDomain       string      `flag:"baseDomain" validate:"required"`
	K8s              *K8sConfig  `default:"{}"`
	Helm             *HelmConfig `default:"{}"`
}

type K8sConfig struct {
	Namespace string `flag:"k8s.namespace"`
}

type HelmConfig struct {
	RepoUrl       string   `flag:"helmRepo" default:"https://raw.githubusercontent.com/ukama/helm-charts/repo-index"`
	Token         string   `flag:"token"`
	IsUpgrade     bool     `flag:"upgrade"`
	ServiceParams []string `flag:"svcParams" default:"[]"`
	SkipDeps      bool     `flag:"skipDeps" default:"false"`
}

type AwsConfig struct {
	AccessKey string `flag:"aws.accessKey"`
	SecretKey string `flag:"aws.secret"`
}
