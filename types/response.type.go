package types

// 기본 응답 구조체
type ResponseType struct {
	Code   string `json:"code"`
	Result bool   `json:"result"`
}

type ResponseSignupType struct {
	Code   string `json:"code"`
	Result bool   `json:"result"`
	BlogId string `json:"blogId"`
}

type ResponseImageUrl struct {
	Code        string   `json:"code"`
	Result      bool     `json:"result"`
	ImageResult []string `json:"imageResult"`
}

// 메세지를 담은 응답
type ResponseMessageType struct {
	Code    string `json:"code"`
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

// 찾은 패스워드 담은 응답
type ResponseFoundPasswordType struct {
	Code     string `json:"code"`
	Result   bool   `json:"result"`
	Password string `json:"password"`
}

// JWT 토큰을 담은 응답
type ResponseTokenType struct {
	Code         string `json:"code"`
	Result       bool   `json:"result"`
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
