package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/iamolegga/enviper"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ukama/ukamaX/cli/pkg"
	"io"
	"os"
	"reflect"
	"strings"
)

// We have to keep all config in one struct.
type GlobalConfig struct {
	Verbose bool
}

/* CodefigReader allows commands to read configuration from config file, env vars or flags.
 Flags have precedence over env vars and evn vars have precedence over config file.
 Flags mapped to config struct filed automatically if their name includes path to field.
 For example:
	type NestedConf struct {
		Foo string
	}

  	type Config struct{
      Nested NestedConf
 	}

in that case flag --ndeste.foo will be mapped automatically
This flag could be also set by UKAMA_NESTED_FOO env var or by creating  config file .ukama.yaml:
nested:
  foo: bar
*/
type ConfigReader interface {
	// ReadConfig reads config from config file, env vars or flags. In case of error fails with os.Exit(1)
	ReadConfig(key string, flags *pflag.FlagSet, rawVal interface{})
	BindFlag(confKey string, flag *pflag.Flag)
}

type ConfMgr struct {
	viper      *enviper.Enviper
	configFile string
	stdout     io.Writer
	stderr     io.Writer
}

func NewConfMgr(configFile string, stdout io.Writer, stderr io.Writer) *ConfMgr {
	return &ConfMgr{
		viper:      enviper.New(viper.New()),
		configFile: configFile,
		stdout:     stdout,
		stderr:     stderr,
	}
}

func (c *ConfMgr) ReadConfig(key string, flags *pflag.FlagSet, rawVal interface{}) {
	c.lateFlagBinding(flags, rawVal)

	if c.configFile != "" {
		// Use config file from the flag.
		c.viper.SetConfigFile(c.configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dmr" (without extension).
		c.viper.AddConfigPath(home)
		c.viper.SetConfigName("." + pkg.CliName)
	}

	c.viper.SetEnvPrefix(strings.ToUpper(pkg.CliName))
	c.viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := c.viper.ReadInConfig(); err == nil {
		fmt.Fprintln(c.stdout, "Using config file:", c.viper.ConfigFileUsed())
	}

	err := c.viper.BindPFlags(flags)
	if err != nil {
		fmt.Fprintf(c.stderr, "Error binding flags. Error: %v", err)
		os.Exit(1)
	}

	err = c.viper.Unmarshal(rawVal)
	if err != nil {
		fmt.Fprintf(c.stderr, "Unable to decode into struct, %v", err)
		os.Exit(1)
	}

	err = c.viper.UnmarshalKey(key, rawVal)
	if err != nil {
		fmt.Fprintf(c.stderr, "Error reading config: '%+v'\n", err)
		os.Exit(1)
	}

	err = validator.New().Struct(rawVal)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if validationErrors != nil && len(validationErrors) > 0 {
			fmt.Fprintf(c.stderr, "Error validating config: '%+v'\n", validationErrors)
			os.Exit(1)
		}
		cobra.CheckErr(err)
	}
}

func (c *ConfMgr) lateFlagBinding(flags *pflag.FlagSet, conf interface{}) {
	t := reflect.TypeOf(conf)
	m := map[string]string{}
	dumpStruct(t, "", m)

	for k, v := range m {
		c.BindFlag(k, flags.Lookup(v))
	}
}

func dumpStruct(t reflect.Type, path string, res map[string]string) {
	switch t.Kind() {
	case reflect.Ptr:
		originalValue := t.Elem()
		dumpStruct(originalValue, path, res)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			val := f.Tag.Get("flag")
			if val != "" && f.Type.Kind() != reflect.Struct && f.Type.Kind() != reflect.Ptr {
				res[strings.TrimPrefix(strings.ToLower(path+"."+f.Name), ".")] = val
			}

			dumpStruct(f.Type, path+"."+f.Name, res)
		}

	case reflect.Interface:
		panic("Interface not supported")

	default:
	}

}

func (c *ConfMgr) BindFlag(confKey string, flag *pflag.Flag) {
	err := c.viper.BindPFlag(confKey, flag)
	if err != nil {
		panic("Unable to bind flag. Error: " + err.Error())
	}
}
