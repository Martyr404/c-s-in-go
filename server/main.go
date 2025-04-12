package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	content := map[string]string{"content": ""}
	router.POST("/upload2", func(ctx *gin.Context) {
		raw_data, err := ctx.GetRawData()
		if err != nil {
			log.Fatal(err)
		} else {
			json.Unmarshal(raw_data, &content)
		}
		if content["content"] != "" {
			passwd := content["content"]
			log.Println(passwd)
		}
	})
	router.Run(":9988")
}
