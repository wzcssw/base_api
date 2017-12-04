package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"tx_base_api/configs"
	"tx_base_api/middleware"
	"tx_base_api/models"
	"tx_base_api/mq"

	"tx_base_api/routers"

	_ "github.com/go-sql-driver/mysql"
	session "github.com/ipfans/echo-session"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// 加载配置文件
	initConfig()
	e := echo.New()
	// e.Debug = true
	// 加载数据库连接
	constr := initDBConfig()
	var err interface{}
	models.Mdb, err = gorm.Open("mysql", constr)
	if err != nil {
		e.Logger.Error(err)
		fmt.Printf("err is :%s \n", err)
		panic("\n数据库连接失败")
	}
	models.Mdb.LogMode(true)
	defer models.Mdb.Close()

	mq.RabbitMQ = mq.InitRabbit()
	mq.RabbitMQ.AddExchangeQueue("tx_base_api")
	mq.RabbitMQ.AddExchangeQueue("txcommon")
	go mq.RabbitMQ.RecivieMsgFromExchangeQueue("tx_base_api", mq.SelfListener)
	go mq.RabbitMQ.RecivieMsgFromExchangeQueue("txcommon", mq.CommonListener)
	defer mq.RabbitMQ.Close()
	// testMsg := mq.MqMsg{
	// 	Ucid:   proto.String("hello"),
	// 	Result: proto.String("msg body"),
	// }
	// mData, err := proto.Marshal(&testMsg)
	// if err != nil {
	// 	fmt.Println("Error1: ", err)
	// 	return
	// }
	// var unMData mq.MqMsg
	// err = proto.Unmarshal(mData, &unMData)
	// if err != nil {
	// 	fmt.Println("Error1: ", err)
	// 	return
	// }

	// 开始其他设置
	e.Use(middleware.CORS())
	// e.Use(middleware.CSRF())
	e.Use(middleware.Secure())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.BodyLimit("2M"))
	e.Static("/static", "static")
	e.Use(middleware.Static("/static"))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 模板引擎设置
	renderer := &mw.TemplateRenderer{
		Templates: template.Must(template.ParseGlob("static/view/*.html")),
	}
	e.Renderer = renderer
	// redis
	redisConfig := configs.AppConfig.Redis
	redisURL := strings.Join([]string{redisConfig.Host, ":", redisConfig.Port}, "")
	store, err := session.NewRedisStore(32, redisConfig.Protocol, redisURL, redisConfig.Password, []byte("tx_base_api"))
	if err != nil {
		panic(err)
	}
	// session
	e.Use(session.Sessions("GSESSION", store))
	// 登录验证
	auth := mw.NewAPIAuth()
	e.Use(auth.CanPass)

	routers.InitRouter(e)

	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Logger.Fatal(e.Start(":4900"))
}

func init() {
	// mq.SendMsgFromRabbit("Hello,World!")
}

func initConfig() {
	// if env == nil {
	// 	env = "development"
	// }
	// 加载配置文件
	// configor.New(&configor.Config{Environment: env}).Load(&AppConfig, "config.yml")
	configor.Load(&configs.AppConfig, "configs/config.yml")
	fmt.Printf("config: %#v \n", configs.AppConfig.DB)
}

// 加载数据库连接
func initDBConfig() string {
	host := configs.AppConfig.DB.Host
	port := configs.AppConfig.DB.Port
	username := configs.AppConfig.DB.Username
	password := configs.AppConfig.DB.Password
	dbname := configs.AppConfig.DB.Database
	charset := configs.AppConfig.DB.Encoding
	constr := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=" + charset + "&parseTime=True&loc=Local"
	fmt.Printf("constr: %s \n", constr)
	return constr
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}
