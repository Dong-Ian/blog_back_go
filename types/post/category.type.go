package types

type ResponseCategoryResponseType struct {
	Code         string   `json:"code"`
	Status       int      `json:"status"`
	Message      string   `json:"message"`
	Result       bool     `json:"result"`
	CategoryList []string `json:"categoryList"`
}

type CategoryQueryResult struct {
	CategoryName string
}
