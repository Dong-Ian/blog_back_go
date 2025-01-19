package response

type CommonResponseWithMessage struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
