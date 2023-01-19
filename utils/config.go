package utils

import "github.com/spf13/viper"

type Config struct {
	InfluxDBToken  string `mapstructure:"INFLUXDB_TOKEN"`
	InfluxDBOrg    string `mapstructure:"INFLUXDB_ORG"`
	InfluxDBBucket string `mapstructure:"INFLUXDB_BUCKET"`
}

// LoadConfig loads config from path
//
// @param path - Path to config file to load
// @param config - Config object with app and env variables
// @param err - Error if there was a problem reading config
//
// @return Config object with env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return
}

func GetEnv(key string) string {
	return viper.GetString(key)
}
