package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/donghquinn/blog_back_go/configs"
	"github.com/donghquinn/blog_back_go/libraries/database"
	"github.com/donghquinn/blog_back_go/types"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 토큰 생성
func CreateJwtToken(userId string, uuid string, userEmail string, userStatus string, blogId string, expire time.Duration) (string, error) {
	globalConfig := configs.GlobalConfig

	redis, redisPingErr := database.RedisInstance()

	if redisPingErr != nil {
		return "", redisPingErr
	}

	getToken, getTokenErr := database.RedisLoginGet(redis, uuid)

	if getTokenErr != nil {
		log.Printf("[JWT] Get Token Error: %v", getTokenErr)
		return "", getTokenErr
	}

	// 이미 등록된 토큰이 있다면 삭제하고 새로 등록
	if getToken != nil {
		log.Printf("[JWT] Found Already Set Token")
		deletErr := database.Delete(redis, userId, uuid)

		if deletErr != nil {
			log.Printf("[JWT] Delete Token Error")
			return "", deletErr
		}
	}

	setLoginErr := database.RedisLoginSet(redis, uuid, userEmail, userStatus, userId, blogId)

	if setLoginErr != nil {
		log.Printf("[JWT] Set Login Error: %v", setLoginErr)

		return "", setLoginErr
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	claims := jwtToken.Claims.(jwt.MapClaims)

	claims["userId"] = userId
	claims["userEmail"] = userEmail
	claims["userStatus"] = userStatus
	claims["uuid"] = uuid
	claims["blogId"] = blogId

	// 만료 시간 - 3시간
	claims["exp"] = time.Now().Add(expire).Unix()

	token, err := jwtToken.SignedString([]byte(globalConfig.JwtKey))

	if err != nil {
		log.Printf("[JWT] Create Token Error: %v", err)

		return "", err
	}

	setErr := database.RedisLoginSet(redis, uuid, userEmail, userStatus, userId, blogId)

	if setErr != nil {
		log.Printf("[JWT] Set Key Error: %v", setErr)
		return "", setErr
	}

	return token, nil
}

// JWT 키  검증
func ValidateJwtTokenFromString(token string) (string, string, string, string, error) {
	redis, redisErr := database.RedisInstance()

	if redisErr != nil {
		return "", "", "", "", redisErr
	}

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

	_, getErr := database.RedisLoginGet(redis, claim.Uuid)

	if getErr != nil {
		log.Printf("[JWT] Get Token Error: %v", getErr)
		return "", "", "", "", getErr
	}

	return claim.UserId, claim.UserEmail, claim.UserStatus, claim.BlogId, nil
}

// JWT 키  검증
func ValidateJwtToken(req *http.Request) (string, string, string, string, error) {
	token := strings.Split(req.Header["Authorization"][0], "Bearer ")[1]
	redis, redisErr := database.RedisInstance()

	if redisErr != nil {
		return "", "", "", "", redisErr
	}

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

	_, getErr := database.RedisLoginGet(redis, claim.Uuid)

	if getErr != nil {
		log.Printf("[JWT] Get Token Error: %v", getErr)
		return "", "", "", "", getErr
	}

	return claim.UserId, claim.UserEmail, claim.UserStatus, claim.BlogId, nil
}
