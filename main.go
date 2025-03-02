package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	defLogName string = "./test.log"
	dsn        string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
)

func logInit() (*os.File, error) {
	file, err := os.OpenFile(defLogName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	return file, nil
}

func storageInit() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	router := gin.Default()

	file, err := logInit()
	if err != nil {
		fmt.Println("Ошибка при инициализации лог-файла")
		return
	}
	defer file.Close()
	db, err := storageInit()
	if err != nil {
		log.Println("Ошибка при инициализации базы данных.")
	}

	router.POST(pingRoute, logMiddleware, ping)
	router.POST(registerRoute, logMiddleware, register(db))
	router.POST(loginRoute, logMiddleware, login(db))
	router.POST(postRoute, logMiddleware, post(db))
	router.POST(commentRoute, logMiddleware, comment(db))

	router.DELETE(delUserRoute, logMiddleware, delUser(db))
	router.DELETE(delPostRoute, logMiddleware, delPost(db))
	router.DELETE(delCommentRoute, logMiddleware, delComment(db))

	router.GET(usersRoute, logMiddleware, allUsers(db))
	router.GET(postsRoute, logMiddleware, userPosts(db))
	router.GET(commentsRoute, logMiddleware, userComments(db))

	router.Run(":8081")
}
