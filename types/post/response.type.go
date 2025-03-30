package types

// 게시글 리스트 응답 구조체
type ResponsePostContentsType struct {
	Status   int                              `json:"status"`
	Code     string                           `json:"code"`
	Message  string                           `json:"message"`
	Result   bool                             `json:"result"`
	PostList ViewSpecificPostContentsResponse `json:"postList"`
}

type ResponsePostRegisterType struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Result  bool   `json:"result"`
	PostSeq int64  `json:"postSeq"`
}

type ResponseGetImageUrlType struct {
	Status      int      `json:"status"`
	Code        string   `json:"code"`
	Message     string   `json:"message"`
	Result      bool     `json:"result"`
	ImageResult []string `json:"imageResult"`
}

// 게시글 리스트 응답 구조체
type ResponsePostListType struct {
	Status    int                         `json:"status"`
	Code      string                      `json:"code"`
	Message   string                      `json:"message"`
	Result    bool                        `json:"result"`
	PostList  []SelectAllPostDataResponse `json:"postList"`
	PostCount int64                       `json:"postCount"`
	Page      string                      `json:"page"`
	Size      string                      `json:"size"`
}

// 게시글 리스트 응답 구조체
type ResponseInsertIdType struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Result   bool   `json:"result"`
	InsertId string `json:"insertId"`
}
