package main

import (
	"fmt"
	"golang_udemy1/40_todo_app/app/controllers"
	"golang_udemy1/40_todo_app/app/models"
)

func main() {
	//特に意味はないが、init()関数を呼び出す為に適当に記述
	fmt.Println(models.Db)

	controllers.StartMainServer()

}
