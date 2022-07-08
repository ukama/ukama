package pkg

const CliName = "ukama"

// We have to keep all config in one struct.
type GlobalConfig struct {
	Verbose bool `flag:"verbose"`
}
