package config

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"testing"
)

const bindedFlag = "id"

type fullConfig struct {
	GlobalConfig `mapstructure:",squash"`
	Conf         LocalConfig
}

type LocalConfig struct {
	Id                 string `flag:"id"`
	FromEnvVar         string
	FromConfig         string
	OverriddenByEvnVar string
	OverriddenByArg    string
}

func Test_ConfigReader(t *testing.T) {
	// arrange
	valFromVar := "valFromVar"
	overridenVar := "overridenVar"
	overrFromArgVal := "overrFromArgValue"

	nc := &fullConfig{}
	confReader := NewConfMgr("testdata/test_conf.yaml")
	cmd := newTestRootCommand(confReader, nc)

	cmd.SetArgs([]string{"get", "--" + bindedFlag, "10", "--verbose", "true", "--conf.overriddenByArg", overrFromArgVal})
	os.Setenv("UKAMA_CONF_FROMENVVAR", valFromVar)
	defer os.Unsetenv("UKAMA_CONF_FROMENVVAR")

	os.Setenv("UKAMA_CONF_OVERRIDDENBYEVNVAR", overridenVar)
	defer os.Unsetenv("UKAMA_CONF_OVERRIDDENBYEVNVAR")

	// act
	err := cmd.Execute()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "10", nc.Conf.Id)
	assert.Equal(t, true, nc.GlobalConfig.Verbose)
	assert.Equal(t, "valFromConf", nc.Conf.FromConfig)
	assert.Equal(t, valFromVar, nc.Conf.FromEnvVar)
	assert.Equal(t, overrFromArgVal, nc.Conf.OverriddenByArg)
}

func newTestRootCommand(confReader ConfigReader, actualConf *fullConfig) *cobra.Command {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "Access node",
	}
	nodeCmd.PersistentFlags().Bool("verbose", false, "verbose")

	nodeCmd.AddCommand(subCommand(confReader, actualConf))
	return nodeCmd
}

// getCmd represents the get command
func subCommand(confReader ConfigReader, actualConf *fullConfig) *cobra.Command {
	getCmd := cobra.Command{
		Use: "get",
		Run: func(cmd *cobra.Command, args []string) {

			err := confReader.ReadConfig("node", cmd.Flags(), actualConf)
			if err != nil {
				log.Fatalf("Failed to read config: %v", err)
			}
		},
	}

	getCmd.Flags().StringP(bindedFlag, "i", "", "id")
	//confReader.BindFlag("conf.id", getCmd.Flags().Lookup(bindedFlag))

	getCmd.Flags().StringP("node.cert", "c", "", "")
	getCmd.Flags().String("conf.overriddenByArg", "", "")

	return &getCmd
}

type dmParent struct {
	GlobalConfig `mapstructure:",squash"`
	Conf         dmSibling `flag:"notAllowed"`
	PtrConf      *dmSibling
	Par          int `flag:"par"`
}

type dmSibling struct {
	Id         string `flag:"id"`
	FromEnvVar string
}

func TestDumpStrunct(t *testing.T) {
	m := map[string]string{}
	dumpStruct(reflect.TypeOf(dmParent{}), "", m)

	assert.Equal(t, "id", m["conf.id"])
	assert.Equal(t, "id", m["ptrconf.id"])
	assert.Equal(t, "par", m["par"])
}
