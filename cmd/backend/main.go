package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/buggoing/echo-template/config"
	"github.com/buggoing/echo-template/database"
	"github.com/buggoing/echo-template/handler"
	"github.com/buggoing/echo-template/logger"
	cusMiddleware "github.com/buggoing/echo-template/middleware"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title pi-cloud-monitor RESTful API
// @version 1.0
// @description The pi-cloud-monitor RESTful API Documents.
// @contact.name wukong
// @contact.email wukong@pplabs.org
// @host 47.114.74.103:20000
// @BasePath /v1

var (
	serverAddr     string
	configFilePath string
	debug          bool
)

var log = logger.New("main")

func init() {
	flag.StringVar(&serverAddr, "addr", ":20000", "the server address")
	flag.StringVar(&configFilePath, "config", "./assets/config.yaml", "the server configuring filepath")
	flag.BoolVar(&debug, "debug", true, "whether debug mode")
	flag.Parse()
	if err := config.Init(configFilePath); err != nil {
		log.Fatalf("fail to load configurations from %s: %v", configFilePath, err)
	}
	log.Warnf("configuration loaded:\n%+v", config.C)
}

func initMysqlDB() (*sql.DB, error) {
	params := map[string]string{
		"charset":   "utf8mb4",
		"collation": "utf8mb4_general_ci",
	}
	m := &mysql.Config{
		User:                 config.C.Mysql.Username,
		Passwd:               config.C.Mysql.Password,
		Net:                  "tcp",
		Addr:                 config.C.Mysql.Address,
		DBName:               config.C.Mysql.DBname,
		AllowNativePasswords: true,
		Params:               params,
	}
	sqlDataSourceName := m.FormatDSN()
	mysqlDB, err := sql.Open("mysql", sqlDataSourceName)
	if err != nil {
		return nil, err
	}
	if err = mysqlDB.Ping(); err != nil {
		return nil, err
	}
	return mysqlDB, nil
}

func mysqlConnMonitor(db *sql.DB) {
	ticker := time.NewTicker(20 * time.Second)
	var err error
	for t := range ticker.C {
		if db == nil {
			if db, err = initMysqlDB(); err != nil {
				log.Errorf("reconnect mysql failed err: %v", err)
			} else {
				log.Warnf("reconnect mysql successfully at %v", t)
			}
			continue
		}
		if err := db.Ping(); err != nil {
			log.Errorf("mysql connection err: %v", err)
			if db, err = initMysqlDB(); err != nil {
				log.Errorf("reconnect mysql failed err: %v", err)
				continue
			}
			log.Warnf("reconnect mysql successfully at %v", t)
		}
	}
}

func main() {
	// influxDBClient := initInfluxDB()
	mysqlClient, err := initMysqlDB()
	if err != nil {
		log.Fatalf("failed to init mysql client err: %v", err)
	}
	go mysqlConnMonitor(mysqlClient)
	dbUser := database.NewUser(mysqlClient)
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	if debug {
		e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			log.Infof("req: %s", reqBody)
			log.Infof("res: %s", resBody)
			log.Info("headers: ")
			for key, value := range c.Request().Header {
				fmt.Printf("%s: %+v\n", key, value)
			}
		}))
	}
	// 不需要 token 认证的 URL
	skippedPath := map[string]bool{
		"/v1/user/login": true,
	}
	authMiddleware := cusMiddleware.JWTMiddleware(config.GetJwtPublicKey(), skippedPath)
	e.GET("/", hello)
	userHandler := handler.NewUser(dbUser)
	v1Group := e.Group("/v1")
	v1Group.Use(authMiddleware)
	userGroup := v1Group.Group("/user")
	userGroup.POST("/login", userHandler.Login)
	e.Logger.Fatal(e.Start(serverAddr))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
