package lconfig

import (
	"context"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"time"
)

type Reader interface {
	SetDefault(key string, value interface{})
	Get(key string) interface{}
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetIntSlice(key string) []int
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	IsSet(key string) bool
	Sub(key string) Reader
	Set(key string, value interface{})
	Save()
}

var configPath *string
var flagset *pflag.FlagSet

type viperConfig struct {
	v      *viper.Viper
	prefix string
	logger zlog.Logger
}

func Init(logger zlog.Logger) (Reader, error) {
	var err error

	flagset = pflag.NewFlagSet("default", pflag.ExitOnError)
	configPath = flagset.String("config", "./unkatan.yml", "configuration file path")

	err = viper.BindPFlags(flagset)
	if err != nil {
		return nil, err
	}

	err = flagset.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	if *configPath == "" {
		return nil, errors.New("empty config file path")
	}

	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigFile(*configPath)
	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if v.GetBool("dump_and_exit") == true {
		fmt.Printf("%+v\n", v.AllSettings())
		os.Exit(0)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logger.Info(
			nil,
			"config file changed",
			zap.String("name", e.Name),
		)
	})

	viperReader := &viperConfig{
		v:      v,
		logger: logger,
	}

	return viperReader, nil
}

func (c *viperConfig) SetDefault(key string, value interface{}) {
	c.v.SetDefault(c.prefix+key, value)
}
func (c *viperConfig) Get(key string) interface{} {
	return c.v.Get(c.prefix + key)
}
func (c *viperConfig) GetBool(key string) bool {
	return c.v.GetBool(c.prefix + key)
}
func (c *viperConfig) GetFloat64(key string) float64 {
	return c.v.GetFloat64(c.prefix + key)
}
func (c *viperConfig) GetInt(key string) int {
	return c.v.GetInt(c.prefix + key)
}
func (c *viperConfig) GetIntSlice(key string) []int {
	return c.v.GetIntSlice(c.prefix + key)
}
func (c *viperConfig) GetString(key string) string {
	return c.v.GetString(c.prefix + key)
}
func (c *viperConfig) GetStringMap(key string) map[string]interface{} {
	return c.v.GetStringMap(c.prefix + key)
}
func (c *viperConfig) GetStringMapString(key string) map[string]string {
	return c.v.GetStringMapString(c.prefix + key)
}
func (c *viperConfig) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(c.prefix + key)
}
func (c *viperConfig) GetTime(key string) time.Time {
	return c.v.GetTime(c.prefix + key)
}
func (c *viperConfig) GetDuration(key string) time.Duration {
	return c.v.GetDuration(c.prefix + key)
}
func (c *viperConfig) IsSet(key string) bool {
	return c.v.IsSet(c.prefix + key)
}
func (c *viperConfig) Sub(key string) Reader {
	var prefix string
	if key == "" {
		prefix = ""
	} else {
		prefix = key + "."
	}
	return &viperConfig{
		v:      c.v,
		prefix: prefix,
		logger: c.logger,
	}
}
func (c viperConfig) Set(key string, value interface{}) {
	c.v.Set(c.prefix+key, value)
}
func (c viperConfig) Save() {
	c.logger.Info(context.Background(), fmt.Sprintf("Saving config to %s", c.v.ConfigFileUsed()))
	err := c.v.WriteConfig()
	if err != nil {
		c.logger.Error(context.Background(), "Config save error", zap.Error(err))
	}
}
