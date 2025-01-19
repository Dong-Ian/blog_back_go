package middlewares

import (
	"net/http"

	"github.com/rs/cors"
)

var originList = []string{
	"http://localhost:3000",
	"https://localhost:3000",
	"https://blog.minjae-dev.com",
	"https://blog.donghyuns.com",
	"unknown",
}

func CorsHanlder() *cors.Cors {
	corHandler := cors.New(cors.Options{
		AllowedOrigins:   originList,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
		Debug:            false,
	})

	return corHandler
}

// func CorsMiddlewares(next *mux.Router) http.Handler {
// 	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
// 		origin := req.Header.Get("Origin")

// 		if origin == "" {
// 			origin = "unknown"
// 		}

// 		for _, o := range originList {
// 			if o == origin {
// 				log.Printf("Allowed Origin: %s", origin)
// 				res.Header().Set("Access-Control-Allow-Origin", origin)
// 				res.Header().Set("Access-Control-Max-Age", "86400")
// 				res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
// 				res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 				res.Header().Set("Access-Control-Allow-Credentials", "true")
// 				break
// 			}
// 		}

// 		res.Header().Set("Access-Control-Allow-Origin", "*")
// 		res.Header().Set("Access-Control-Max-Age", "86400")
// 		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
// 		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		res.Header().Set("Access-Control-Allow-Credentials", "true")

// 		// Handle preflight request
// 		if req.Method == http.MethodOptions {
// 			res.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(res, req)
// 	})
// }
