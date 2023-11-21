package utils

import (
	"os"

	"github.com/num30/config"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/testing/integration/pkg"

	log "github.com/sirupsen/logrus"
)

func LoadConfigVarsFromEnv(serviceName string, serviceConfig *pkg.Config) {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)

	err := config.NewConfReader(serviceName).Read(serviceConfig)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if serviceConfig.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
}
