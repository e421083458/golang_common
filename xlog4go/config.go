package xlog4go

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

func SetupLogWithConf(file string) (err error) {
	var lc LogConfig
	cnt, err := ioutil.ReadFile(file)

	//var conf Config
	if _, err := toml.Decode(string(cnt), &lc); err != nil {
		return err
	}

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
			Register(w)
		}

		if len(lc.FW.WfLogPath) > 0 {
			wfw := NewFileWriter()
			wfw.SetFileName(lc.FW.WfLogPath)
			wfw.SetPathPattern(lc.FW.RotateWfLogPath)
			wfw.SetLogLevelFloor(WARNING)
			wfw.SetLogLevelCeil(ERROR)
			Register(wfw)
		}
	}

	if lc.CW.On {
		w := NewConsoleWriter()
		w.SetColor(lc.CW.Color)
		Register(w)
	}

	switch lc.Level {
	case "trace":
		SetLevel(TRACE)

	case "debug":
		SetLevel(DEBUG)

	case "info":
		SetLevel(INFO)

	case "warning":
		SetLevel(WARNING)

	case "error":
		SetLevel(ERROR)

	case "fatal":
		SetLevel(FATAL)

	default:
		err = errors.New("Invalid log level")
	}
	return
}
