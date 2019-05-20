package log

import (
	"errors"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type ConfFileWriter struct {
	On              bool   `toml:"On"`
	LogPath         string `toml:"LogPath"`
	RotateLogPath   string `toml:"RotateLogPath"`
	WfLogPath       string `toml:"WfLogPath"`
	RotateWfLogPath string `toml:"RotateWfLogPath"`
}

type ConfConsoleWriter struct {
	On    bool `toml:"On"`
	Color bool `toml:"Color"`
}

type LogConfig struct {
	Level string            `toml:"LogLevel"`
	FW    ConfFileWriter    `toml:"FileWriter"`
	CW    ConfConsoleWriter `toml:"ConsoleWriter"`
}

func SetupLogInstanceWithConf(lc LogConfig,logger *Logger) (err error) {
	if lc.FW.On {
		if len(lc.FW.LogPath) > 0 {
			w := NewFileWriter()
			w.SetFileName(lc.FW.LogPath)
			w.SetPathPattern(lc.FW.RotateLogPath)
			w.SetLogLevelFloor(TRACE)
			if len(lc.FW.WfLogPath) > 0 {
				w.SetLogLevelCeil(INFO)
			} else {
				w.SetLogLevelCeil(ERROR)
			}
			logger.Register(w)
		}

		if len(lc.FW.WfLogPath) > 0 {
			wfw := NewFileWriter()
			wfw.SetFileName(lc.FW.WfLogPath)
			wfw.SetPathPattern(lc.FW.RotateWfLogPath)
			wfw.SetLogLevelFloor(WARNING)
			wfw.SetLogLevelCeil(ERROR)
			logger.Register(wfw)
		}
	}

	if lc.CW.On {
		w := NewConsoleWriter()
		w.SetColor(lc.CW.Color)
		logger.Register(w)
	}
	switch lc.Level {
	case "trace":
		logger.SetLevel(TRACE)

	case "debug":
		logger.SetLevel(DEBUG)

	case "info":
		logger.SetLevel(INFO)

	case "warning":
		logger.SetLevel(WARNING)

	case "error":
		logger.SetLevel(ERROR)

	case "fatal":
		logger.SetLevel(FATAL)

	default:
		err = errors.New("Invalid log level")
	}
	return
}

func SetupLogInstanceWithFile(file string,logger *Logger) (err error) {
	var lc LogConfig
	cnt, err := ioutil.ReadFile(file)
	if _, err := toml.Decode(string(cnt), &lc); err != nil {
		return err
	}
	return SetupLogInstanceWithConf(lc, logger)
}

func SetupDefaultLogWithFile(file string) (err error) {
	defaultLoggerInit()
	return SetupLogInstanceWithFile(file, logger_default)
}

func SetupDefaultLogWithConf(lc LogConfig) (err error) {
	defaultLoggerInit()
	return SetupLogInstanceWithConf(lc, logger_default)
}
