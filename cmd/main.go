package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dingkegithub/user/dao"
	"github.com/dingkegithub/user/endpoint"
	"github.com/dingkegithub/user/redis"
	"github.com/dingkegithub/user/service"
	"github.com/dingkegithub/user/transport"
)

func main() {
	var (
		servicePort = flag.Int("service.port", 10086, "service port")

		mysqlAddr = flag.String("mysql.addr", "127.0.0.1", "mysql addr")

		mysqlPort = flag.String("mysql.port", "3306", "mysql port")

		redisAddr = flag.String("redis.addr", "127.0.0.1", "redis addr")

		redisPort = flag.String("redis.port", "6379", "redis port")
	)

	flag.Parse()

	ctx := context.Background()

	if err := dao.MysqlInit(*mysqlAddr, *mysqlPort, "root", "123456", "user"); err != nil {
		log.Fatal(err)
	}

	if err := redis.InitRedis(*redisAddr, *redisPort, ""); err != nil {
		log.Fatal(err)
	}

	userService := service.MakeUserService(&dao.UserDaoImpl{})
	userEndpoints := endpoint.NewUserEndpoints(userService)

	router := transport.MakeHttpHandler(ctx, userEndpoints)

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	log.Println(<-errChan)
}
