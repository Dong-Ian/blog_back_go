package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/donghquinn/blog_back_go/auth"
	"github.com/donghquinn/blog_back_go/response"
)

// 사용자 정의 키 타입을 사용하여 컨텍스트 충돌 방지
type contextKey string

const (
	// JWT 서명에 사용할 비밀 키 (환경 변수로 관리하는 것이 좋습니다)
	// jwtSecret = configs.GlobalConfig.JwtKey // 실제 배포 시 환경 변수로 관리

	// 컨텍스트에 사용자 정보를 저장할 키
	userContextKey = contextKey("user")
)

// 사용자 정보 구조체
type User struct {
	UserId     string
	UserEmail  string
	UserStatus string
	BlogId     string
}

var excludeRouteList = []string{
	"/", "/api", "/refresh",
	"/user/signup", "/user/login", "/user/search/email", "/user/search/password", "/user/profile",
	"/post/contents", "/post/list", "/post/list/pinned", "/post/list/tag", "/post/list/category", "/post/url", "/post/category/list",
}

// AuthMiddleware는 accessToken 쿠키를 추출하고 JWT를 검증하는 미들웨어입니다.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 제외할 경로는 바로 다음 핸들러로 넘김
		for _, route := range excludeRouteList {
			if strings.HasPrefix(r.URL.Path, route) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// accessToken 쿠키 추출
		cookie, err := r.Cookie("accessToken")
		if err != nil {
			if err == http.ErrNoCookie {
				response.Response(w, response.CommonResponseWithMessage{
					Status:  http.StatusUnauthorized,
					Code:    "AUTH001",
					Message: "No access token provided",
				})
				return
			}
			// 다른 쿠키 에러 처리
			response.Response(w, response.CommonResponseWithMessage{
				Status:  http.StatusUnauthorized,
				Code:    "AUTH002",
				Message: "Invalid cookie format",
			})
			return
		}

		accessToken := cookie.Value

		userId, userEmail, userStatus, blogId, validateErr := auth.ValidateJwtTokenFromString(accessToken)

		if validateErr != nil {
			// 토큰 만료에 대한 응답
			if strings.Contains(validateErr.Error(), "token expired") {
				response.Response(w, response.CommonResponseWithMessage{
					Status:  http.StatusUnauthorized,
					Code:    "AUTH003",
					Message: "Token expired",
				})
				return
			}

			// 일반적인 JWT 검증 실패 응답
			response.Response(w, response.CommonResponseWithMessage{
				Status:  http.StatusUnauthorized,
				Code:    "AUTH004",
				Message: "Invalid token",
			})
			return
		}
		// 사용자 정보를 구조체로 생성
		user := User{
			UserId:     userId,
			UserEmail:  userEmail,
			UserStatus: userStatus,
			BlogId:     blogId,
		}

		// 사용자 정보를 컨텍스트에 추가
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 사용자 정보를 핸들러에서 가져오는 헬퍼 함수
func GetUserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}
