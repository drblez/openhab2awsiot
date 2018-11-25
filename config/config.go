package config

import (
	"github.com/jessevdk/go-flags"
	"github.com/joomcode/errorx"
)

var (
	Errors       = errorx.NewNamespace("config")
	CommonErrors = Errors.NewType("common_error")
	UnknownFlag  = Errors.NewType("unknown_flags")
)

type Config struct {
	RunAsProgram bool `short:"R" description:"Run as program"`
	Debug        bool `long:"debug" description:"Debug level logging" env:"DEBUG"`
	Console      bool `long:"console" description:"Output to console" env:"CONSOLE"`
}

func InitConfig() (*Config, error) {
	config := &Config{}
	f := flags.NewParser(config, flags.Default)
	_, err := f.Parse()
	if err != nil {
		switch err := err.(type) {
		case *flags.Error:
			switch err.Type {
			case flags.ErrUnknownFlag:
				return nil, UnknownFlag.New(err.Message)
			}
		}
		return nil, CommonErrors.Wrap(err, "Config parse error")
	}
	return config, nil
}
