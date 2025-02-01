package main

import (
	"fmt"
	"log"
	"os"
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

func main() {
	// router := gin.Default()

	file, err := logInit()
	if err != nil {
		fmt.Println("Ошибка при инициализации лог-файла")
		return
	}
	defer file.Close()
}
