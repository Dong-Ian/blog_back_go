package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/donghquinn/blog_back_go/configs"
	"github.com/donghquinn/blog_back_go/types"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 토큰 생성
func CreateJwtToken(userId string, userEmail string, userStatus string, blogId string, expire time.Duration) (string, error) {
	globalConfig := configs.GlobalConfig

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	claims := jwtToken.Claims.(jwt.MapClaims)

	claims["userId"] = userId
	claims["userEmail"] = userEmail
	claims["userStatus"] = userStatus
	claims["blogId"] = blogId

	// 만료 시간 - 3시간
	claims["exp"] = time.Now().Add(expire).Unix()

	token, err := jwtToken.SignedString([]byte(globalConfig.JwtKey))

	if err != nil {
		log.Printf("[JWT] Create Token Error: %v", err)

		return "", err
	}

	return token, nil
}

// JWT 키  검증
func ValidateJwtTokenFromString(token string) (string, string, string, string, error) {
	globalConfig := configs.GlobalConfig

	// JWT 토큰 파싱
	parseToken, err := jwt.ParseWithClaims(token, &types.JwtInterface{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			parseErr := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

			log.Printf("[JWT] Parse With Claims Error: %v", parseErr)
			return nil, parseErr
		}

		return []byte(globalConfig.JwtKey), nil
	})

	if err != nil {
		log.Printf("[JWT] Parsing JWT Validation Error: %v", err)

		return "", "", "", "", err
	}

	claim, ok := parseToken.Claims.(*types.JwtInterface)

	if !ok {
		claimErr := fmt.Errorf("can't parse values from token")
		log.Printf("[JWT] Parse Token with Claims: %v", claimErr)
		return "", "", "", "", claimErr
	}

	return claim.UserId, claim.UserEmail, claim.UserStatus, claim.BlogId, nil
}

// JWT 키  검증
func ValidateJwtToken(req *http.Request) (string, string, string, string, error) {
	token := strings.Split(req.Header["Authorization"][0], "Bearer ")[1]

	globalConfig := configs.GlobalConfig

	// JWT 토큰 파싱
	parseToken, err := jwt.ParseWithClaims(token, &types.JwtInterface{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			parseErr := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

			log.Printf("[JWT] Parse With Claims Error: %v", parseErr)
			return nil, parseErr
		}

		return []byte(globalConfig.JwtKey), nil
	})

	if err != nil {
		log.Printf("[JWT] Parsing JWT Validation Error: %v", err)

		return "", "", "", "", err
	}

	claim, ok := parseToken.Claims.(*types.JwtInterface)

	if !ok {
		claimErr := fmt.Errorf("can't parse values from token")
		log.Printf("[JWT] Parse Token with Claims: %v", claimErr)
		return "", "", "", "", claimErr
	}

	return claim.UserId, claim.UserEmail, claim.UserStatus, claim.BlogId, nil
}
