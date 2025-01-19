package middlewares

import (
	"log"
	"net/http"
)

var originList = []string{
	"http://localhost:3000",
	"https://blog.minjae-dev.com/",
	"https://blog.donghyuns.com",
	"unknown",
}

func CorsMiddlewares(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Println("CORS middleware reached") // 호출 확인

		origin := req.Header.Get("Origin")
		log.Printf("Origin: %s", origin)

		// Origin 헤더가 없으면 기본 설정
		if origin == "" {
			origin = "unknown"
		}

		// 요청의 Origin이 허용된 Origin 목록에 있는지 확인
		isAllowed := false
		for _, allowedOrigin := range originList {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			res.Header().Set("Access-Control-Allow-Origin", origin)
			res.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			log.Printf("Origin not allowed: %s", origin)
		}

		res.Header().Set("Access-Control-Max-Age", "86400")
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if req.Method == http.MethodOptions {
			log.Println("Handling OPTIONS request")
			res.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(res, req)
	})
}
