package configs

// AppConfig 全局配置文件
var AppConfig = struct {
	APPName string `default:"tx_base_api"`

	DB struct {
		Host     string `default:"localhost"`
		Pool     string `default:"5"`
		Database string `default:"tx_base_api"`
		Username string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     string `default:"3306"`
		Encoding string `default:"utf-8"`
	}

	Redis struct {
		Host     string `default:"localhost"`
		Port     string `default:"6379"`
		Password string `default:""`
		Protocol string `default:"tcp"`
	}

	Rabbit struct {
		URL string `default:"amqp://rabbitmq:123456@localhost:5672/"`
	}
}{}
