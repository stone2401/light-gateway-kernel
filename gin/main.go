package main

import "github.com/gin-gonic/gin"

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *Login) UnmarshalJSON(data []byte) error {
	return nil
}

func main() {
	app := gin.Default()
	app.GET("/hello", func(c *gin.Context) {
		login := &Login{}
		c.ShouldBindJSON(login)
	})
	app.Run(":8080")
}
