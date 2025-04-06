package network

import (
	"log"
	"net/http"
	"time"

	"github.com/donghquinn/blog_back_go/configs"
	"github.com/donghquinn/blog_back_go/libraries/database"
	"github.com/donghquinn/blog_back_go/middlewares"
	"github.com/donghquinn/blog_back_go/routers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func OpenServer() *http.Server {
	router := mux.NewRouter()

	router.Use(middlewares.CheckBlogId)

	routers.DefaultRouter(router)
	routers.UploadImageRouter(router)

	routers.PostAdminRouter(router)
	routers.PostRouter(router)

	routers.UserAdminRouter(router)
	routers.UserRouter(router)

	// handler := middlewares.CorsMiddlewares(router)
	// handler := cors.Default().Handler(router)

	handler := middlewares.CorsHanlder().Handler(router)
	// router.Use(mux.CORSMethodMiddleware(router))

	serving := &http.Server{
		Handler:      handler,
		Addr:         configs.GlobalConfig.AppHost,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	return serving
}

func DatabaseConnect() {
	minioErr := database.MinioConnect()

	if minioErr != nil {
		log.Printf("[START] Minio Connection Check Error: %v", minioErr)
	}

	checkErr := database.CheckConnection()

	if checkErr != nil {
		log.Printf("[START] Databae Connection Check Error: %v", checkErr)
	}

	// _, redisErr := database.RedisInstance()

	// if redisErr != nil {
	// 	log.Printf("[START] Redis Connection Check Error: %v", redisErr)
	// }

}

func SetConfigs() {
	envErr := godotenv.Load(".env")

	if envErr != nil {
		log.Printf("[ENV] Load Env Error")
	}

	configs.SetGlobalConfig()
	configs.SetDatabaseConfig()
	configs.SetMinioConfig()
	configs.SetRedisConfig()
}
