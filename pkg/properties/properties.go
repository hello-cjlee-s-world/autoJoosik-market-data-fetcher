package properties

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
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
	Port            int    `mapstructure:"port"`
	Name            string `mapstructure:"name"`
	MaxIdleConn     int    `mapstructure:"maxIdleConn"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
	SSLMode         string `mapstructure:"sslMode"`
	MaximumPoolSize int    `mapstructure:"maximumPoolSize"`
}

type kiwoomApiStruct struct {
	AppKey    string `mapstructure:"appKey"`
	SecretKey string `mapstructure:"secretKey"`
}

//type buyConstraintsStruct struct {
//MaxHoldingCount      int `mapstructure:"maxHoldingCount"`
//MaxDailyBuyCount     int `mapstructure:"maxDailyBuyCount"`
//CooldownAfterBuy     int `mapstructure:"cooldownAfterBuy"`
//AllowAddBuy          int `mapstructure:"allowAddBuy"`
//MaxInvestPerStockPct int `mapstructure:"maxInvestPerStockPct"`
//}

type PropertiesInfo struct {
	configFile string
	Server     serverStruct    `mapstructure:"server"`
	Logging    loggingStruct   `mapstructure:"logging"`
	Database   databaseStruct  `mapstructure:"database"`
	KiwoomApi  kiwoomApiStruct `mapstructure:"kiwoomApi"`
	//BuyConstraints buyConstraintsStruct `mapstructure:"buyConstraints"`
}

func (p *PropertiesInfo) Init(configFile string) {
	p.configFile = configFile

	viper.SetConfigFile(configFile) // name of config file (without extension)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
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
