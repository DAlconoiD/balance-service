package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"strings"
)

//Config - application config
type Config struct {
	ServerAddress      string
	DBConnectionString string
	PaginationNumber   int
}

//LoadConfig loads config from path=p
func LoadConfig(p string) (*Config, error) {
	dir := path.Dir(p)
	_, file := path.Split(p)
	parts := strings.Split(file, ".")
	filename := parts[0]
	extension := parts[1]

	viper.AddConfigPath(dir)
	viper.SetConfigName(filename)
	viper.SetConfigType(extension)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := Config{}

	viper.SetDefault("DB.HOST", "localhost")
	viper.SetDefault("DB.USER", "postgres")
	viper.SetDefault("DB.PASSWORD", "password")
	viper.SetDefault("DB.NAME", "balance")
	viper.SetDefault("DB.PORT", "5432")
	viper.SetDefault("DB.SSL", "disable")
	dbHost := viper.Get("DB.HOST")
	dbUser := viper.Get("DB.USER")
	dbPwd := viper.Get("DB.PASSWORD")
	dbName := viper.Get("DB.NAME")
	dbPort := viper.Get("DB.PORT")
	dbSSL := viper.Get("DB.SSL")
	config.DBConnectionString = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", dbHost, dbUser, dbPwd, dbName, dbPort, dbSSL)

	viper.SetDefault("SERVER.PORT", "8080")
	srvPort := viper.Get("SERVER.PORT")
	config.ServerAddress = fmt.Sprintf(":%v", srvPort)

	viper.SetDefault("SETTINGS.PAGINATION_NUM", 10)
	config.PaginationNumber = viper.GetInt("SETTINGS.PAGINATION_NUM")

	return &config, nil
}
