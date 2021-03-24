OS=linux
ARCH=amd64

target=user
src=cmd/main.go

${target}: pre ${src}
	GOOS=linux GOARCH=amd64 go build -o $@ ${src}

pre:
	mkdir -p bin

.PHONY:clean
clean:
	rm -rf bin

.PHONY:mysql
mysql:
	docker build -t mysql-for-user -f Dockerfiles/MysqlDockerfile .
	docker run -itd --name mysql-for-user -p 23306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql-for-user

.PHONY:redis
redis:
	docker run -itd --name redis-user -p 26379:6379 redis:latest

.PHONY:usersrv
usersrv:
	docker build -t user -f Dockerfiles/UserDockerfile .
