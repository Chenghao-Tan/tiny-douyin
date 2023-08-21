package conf

type Redis struct {
	RedisHost string `yaml:"redisHost"`
	RedisPort string `yaml:"redisPort"`
	RedisDB   int    `yaml:"redisDB"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	TLS       bool   `yaml:"tls"`
}
