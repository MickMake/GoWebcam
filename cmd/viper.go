package cmd


// var rootViper *viper.Viper
//
// // initConfig reads in config file and ENV variables if set.
// func initConfig(cmd *cobra.Command) error {
// 	var err error
//
// 	for range Only.Once {
// 		rootViper = viper.New()
// 		rootViper.AddConfigPath(Cmd.ConfigDir)
// 		rootViper.SetConfigFile(Cmd.ConfigFile)
// 		// rootViper.SetConfigName("config")
//
// 		// If a config file is found, read it in.
// 		err = openConfig()
// 		if err != nil {
// 			break
// 		}
//
// 		rootViper.SetEnvPrefix(defaults.EnvPrefix)
// 		rootViper.AutomaticEnv() // read in environment variables that match
// 		err = bindFlags(cmd, rootViper)
// 		if err != nil {
// 			break
// 		}
// 	}
//
// 	return err
// }
//
// func openConfig() error {
// 	var err error
//
// 	for range Only.Once {
// 		err = rootViper.ReadInConfig()
// 		if _, ok := err.(viper.UnsupportedConfigError); ok {
// 			break
// 		}
//
// 		if _, ok := err.(viper.ConfigParseError); ok {
// 			break
// 		}
//
// 		if _, ok := err.(viper.ConfigMarshalError); ok {
// 			break
// 		}
//
// 		if os.IsNotExist(err) {
// 			rootViper.SetDefault(flagDebug, Cmd.Debug)
// 			rootViper.SetDefault(flagQuiet, Cmd.Quiet)
//
// 			rootViper.SetDefault(flagWebHost, Cmd.WebHost)
// 			rootViper.SetDefault(flagWebPort, Cmd.WebPort)
// 			rootViper.SetDefault(flagWebUsername, Cmd.WebUsername)
// 			rootViper.SetDefault(flagWebPassword, Cmd.WebPassword)
// 			rootViper.SetDefault(flagWebTimeout, Cmd.WebTimeout)
//
// 			err = rootViper.WriteConfig()
// 			if err != nil {
// 				break
// 			}
//
// 			err = rootViper.ReadInConfig()
// 			break
// 		}
// 		if err != nil {
// 			break
// 		}
//
// 		err = rootViper.MergeInConfig()
// 		if err != nil {
// 			break
// 		}
//
// 		// err = viper.Unmarshal(Cmd)
// 	}
//
// 	return err
// }
//
// func writeConfig() error {
// 	var err error
//
// 	for range Only.Once {
// 		err = rootViper.MergeInConfig()
// 		if err != nil {
// 			break
// 		}
//
// 		rootViper.Set(flagDebug, Cmd.Debug)
// 		rootViper.Set(flagQuiet, Cmd.Quiet)
//
// 		rootViper.Set(flagWebHost, Cmd.WebHost)
// 		rootViper.Set(flagWebPort, Cmd.WebPort)
// 		rootViper.Set(flagWebUsername, Cmd.WebUsername)
// 		rootViper.Set(flagWebPassword, Cmd.WebPassword)
// 		rootViper.Set(flagWebTimeout, Cmd.WebTimeout)
//
// 		err = rootViper.WriteConfig()
// 		if err != nil {
// 			break
// 		}
// 	}
//
// 	return err
// }
//
// func readConfig() error {
// 	var err error
//
// 	for range Only.Once {
// 		err = rootViper.ReadInConfig()
// 		if err != nil {
// 			break
// 		}
//
// 		_, _ = fmt.Fprintln(os.Stderr, "Config file settings:")
//
// 		_, _ = fmt.Fprintf(os.Stderr, "Web Host:			%v\n", rootViper.Get(flagWebHost))
// 		_, _ = fmt.Fprintf(os.Stderr, "Web Port:			%v\n", rootViper.Get(flagWebPort))
// 		_, _ = fmt.Fprintf(os.Stderr, "Web UserAccount:	%v\n", rootViper.Get(flagWebUsername))
// 		_, _ = fmt.Fprintf(os.Stderr, "Web UserPassword:	%v\n", rootViper.Get(flagWebPassword))
// 		_, _ = fmt.Fprintf(os.Stderr, "Web Timeout:		%v\n", rootViper.Get(flagWebPort))
// 		_, _ = fmt.Fprintln(os.Stderr)
//
// 		_, _ = fmt.Fprintf(os.Stderr, "Debug:		%v\n", rootViper.Get(flagDebug))
// 		_, _ = fmt.Fprintf(os.Stderr, "Quiet:		%v\n", rootViper.Get(flagQuiet))
// 	}
//
// 	return err
// }
//
// func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
// 	var err error
//
// 	cmd.Flags().VisitAll(func(f *pflag.Flag) {
// 		// Environment variables can't have dashes in them, so bind them to their equivalent
// 		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
// 		if strings.Contains(f.Name, "-") {
// 			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
// 			err = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", defaults.EnvPrefix, envVarSuffix))
// 		}
//
// 		// Apply the viper config value to the flag when the flag is not set and viper has a value
// 		if !f.Changed && v.IsSet(f.Name) {
// 			val := v.Get(f.Name)
// 			err = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
// 		}
// 	})
//
// 	return err
// }
