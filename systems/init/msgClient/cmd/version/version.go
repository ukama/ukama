package version

// Minor version is autoupdated by the build system
<<<<<<<< HEAD:systems/subscriber/api-gateway/cmd/version/version.go
// NOTE: use go build -ldflags "-X github.com/ukama/ukama/systems/.../cmd/version.Version==$(git describe)"
========
// NOTE: use go build -ldflags "-X github.com/ukama/ukama/systems/init/msgClient/cmd/version.Version==$(git describe)"
>>>>>>>> origin/subscriber-system:systems/init/msgClient/cmd/version/version.go
var Version = "v0.0.debug"
