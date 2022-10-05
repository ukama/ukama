package version

// Minor version is autoupdated by the build system
// NOTE: use go build -ldflags "-X github.com/ukama/ukama/systems/registry/users/cmd/version.Version==$(git describe)"
var Version = "v0.0.debug"
