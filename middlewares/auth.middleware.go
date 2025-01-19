package middlewares

import (
	"context"
	"net/http"

	"github.com/donghquinn/blog_back_go/auth"
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
	UserId    string
	UserEmail string
	UserType  string
	BlogId    string
}

// AuthMiddleware는 accessToken 쿠키를 추출하고 JWT를 검증하는 미들웨어입니다.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// accessToken 쿠키 추출
		cookie, err := r.Cookie("accessToken")

		if err != nil {
			if err == http.ErrNoCookie {
				// 쿠키가 없을 경우 401 Unauthorized 응답
				http.Error(w, "Unauthorized: No access token", http.StatusUnauthorized)
				return
			}
			// 다른 쿠키 에러 처리
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		accessToken := cookie.Value

		userId, userEmail, userType, blogId, validateErr := auth.ValidateJwtTokenFromString(accessToken)

		if validateErr != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// 사용자 정보를 구조체로 생성
		user := User{
			UserId:    userId,
			UserEmail: userEmail,
			UserType:  userType,
			BlogId:    blogId,
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
