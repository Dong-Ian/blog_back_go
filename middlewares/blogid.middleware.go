package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/donghquinn/blog_back_go/response"
)

// CheckBlogId는 요청 헤더에 BlogId가 존재하는지 확인하는 미들웨어입니다.
func CheckBlogId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		blogId := req.Header.Get("BlogId")

		if blogId == "" {
			log.Printf("[GET_POST] BlogId header is missing")
			response.Response(res, response.CommonResponseWithMessage{
				Status:  http.StatusBadRequest,
				Code:    "01",
				Message: "No blog Id Provided",
			})
			return
		}

		// BlogId를 컨텍스트에 저장
		ctx := context.WithValue(req.Context(), "BlogId", blogId)

		// 새 컨텍스트로 요청 업데이트
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// 여러 미들웨어를 체인으로 연결할 경우 사용할 수 있는 헬퍼 함수
func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
