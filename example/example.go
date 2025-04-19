package main

import (
	"log"

	"github.com/HARB1I/cin"
)

func main() {
	router := cin.New()

	router.POST("/auth", Auth)
	auth := router.Group("/auth", AuthMiddleWare)
	{
		auth.GET("/get/{id}", GetTest)
	}

	log.Fatal(router.Run(":8080"))
}

type User struct {
	Name string `json:"name"`
}

type Error struct {
	Error string `json:"error"`
}

func GetTest(c *cin.Context) cin.Response {
	id := c.PathValue("id")
	println(id)

	return cin.Resp(cin.StatusOK, cin.H{
		"msg": "hello",
	})
}

func Auth(c *cin.Context) cin.Response {
	user := User{}
	err := c.BindJSON(&user)
	if err != nil {
		println(err)
		return cin.Resp(cin.StatusBadRequest, Error{
			Error: "error bad request",
		})
	}
	return cin.Resp(cin.StatusCreated, user)
}

func AuthMiddleWare(c *cin.Context) cin.Response {
	// return cin.Resp(code, obj) если не хотите пропускать

	return nil // если хотите пропустить
}
