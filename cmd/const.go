package cmd

import "time"

//goland:noinspection SpellCheckingInspection
const (
	EnvPrefix         = "WEBCAM"
	defaultConfigFile = "config.json"

	flagConfigFile = "config"
	flagDebug      = "debug"
	flagQuiet      = "quiet"
	flagDaemonize   = "daemon"

	flagWebHost     = "host"
	flagWebPort     = "port"
	flagWebUsername = "user"
	flagWebPassword = "password"
	flagWebTimeout  = "timeout"
	flagWebPrefix   = "prefix"

	defaultHost     = ""
	defaultPort     = ""
	defaultUsername = ""
	defaultPassword = ""
	defaultTimeout  = time.Second * 30
	defaultPrefix   = ""
)
