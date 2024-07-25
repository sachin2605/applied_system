package main

import (
	"applied_system/server"
	"fmt"
)

func main() {
	router := server.SetupRouter()
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("failed to run server: %v\n", err)
	}
}
