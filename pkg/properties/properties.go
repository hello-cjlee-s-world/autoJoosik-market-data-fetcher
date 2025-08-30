package properties

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type serverStruct struct {
	Port        int    `mapstructure:"port"`
	Title       string `mapstructure:"title"`
	ContextPath string `mapstructure:"contextPath"`
}

type loggingStruct struct {
	Level         string `mapstructure:"level"`
	Filename      string `mapstructure:"filename"`
	MaxSize       int    `mapstructure:"maxSize"`
	MaxBackups    int    `mapstructure:"maxBackups"`
	MaxAge        int    `mapstructure:"maxAge"`
	Compress      bool   `mapstructure:"compress"`
	ConsoleOutput bool   `mapstructure:"consoleOutput"`
}

type databaseStruct struct {
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Name            string `mapstructure:"name"`
	MaxIdleConn     int    `mapstructure:"maxIdleConn"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
}

type PropertiesInfo struct {
	configFile string
	Server     serverStruct   `mapstructure:"server"`
	Logging    loggingStruct  `mapstructure:"logging"`
	Database   databaseStruct `mapstructure:"database"`
}

func (p *PropertiesInfo) Init(configFile string) {
	p.configFile = configFile

	viper.SetConfigFile(configFile) // name of config file (without extension)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		fmt.Println("Properties init error:", err)
		logger.Error("Properties init",
			"msg", err.Error(),
		)
		panic(err)
	}

	err1 := viper.Unmarshal(p)
	if err1 != nil {
		logger.Error("Properties init",
			"msg", err1.Error(),
		)
		panic(err1)
	}
}

func (p *PropertiesInfo) GetString(key string) string { return viper.GetString(key) }

var (
	instance *PropertiesInfo
	once     sync.Once
)

func GetInstance() *PropertiesInfo {
	once.Do(func() {
		instance = &PropertiesInfo{}
	})
	return instance
}
