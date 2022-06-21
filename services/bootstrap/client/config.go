package bootstrap

type AuthConfig struct {
	ClientId     string
	ClientSecret string
	Audience     string `default:"bootstrap.ukama.com"`
	GrantType    string `default:"client_credentials"`
	Auth0Host    string
}

type DebugConf struct {
	DisableBootstrap bool
}
