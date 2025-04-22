package main

import (
	"annisa-api/handler"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	handler.StartApp()
}
