package main

import (
	"log"

	"github.com/donghquinn/blog_back_go/configs"
	"github.com/donghquinn/blog_back_go/libraries/network"
	testlogic "github.com/donghquinn/blog_back_go/test"
)

func main() {
	network.SetConfigs()

	network.DatabaseConnect()
	
	serving := network.OpenServer()

	testlogic.Test()
	
	log.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	log.Printf("[DEBUG] App Host %s", configs.GlobalConfig.AppHost)
	log.Printf("[START] Server Listening On: %s", configs.GlobalConfig.AppPort)
	log.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")


	serving.ListenAndServe()
}

