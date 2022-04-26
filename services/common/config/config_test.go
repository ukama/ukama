package config

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	BaseConfig `mapstructure:",squash"`
	SomeUrl    string
	DB         *Database // this should be initialized
}

func TestLoadConfig(t *testing.T) {
	const url = "test_url"

	t.Run("EnvVars", func(t *testing.T) {
		os.Setenv("SOMEURL", url)
		os.Setenv("DB_DBNAME", "connStr")
		os.Setenv("DEBUGMODE", "true")
		conf := &TestConfig{
			DB: &Database{},
		}
		LoadConfig("test-config", conf)
		assert.Equal(t, url, conf.SomeUrl)
		assert.Equal(t, "connStr", conf.DB.DbName)
		assert.Equal(t, true, conf.DebugMode)

		os.Unsetenv("SOMEURL")
		os.Unsetenv("DB_DBNAME")
		os.Unsetenv("DEBUGMODE")
	})

	t.Run("ConfigFile", func(t *testing.T) {
		home, err := homedir.Dir()
		assert.NoError(t, err)

		file := path.Join(home, "test-config.yaml")
		fileContent := `
someUrl: test_url
debugMode: "true"
db:
   dbName: connectionStr
`
		err = ioutil.WriteFile(file, []byte(fileContent), 0644)
		assert.NoError(t, err)

		conf := &TestConfig{}
		LoadConfig("test-config", conf)
		assert.Equal(t, url, conf.SomeUrl)
		assert.Equal(t, "connectionStr", conf.DB.DbName)
		assert.Equal(t, true, conf.DebugMode)
		os.Remove(file)
	})

}
