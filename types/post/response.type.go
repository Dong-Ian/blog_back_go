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

type ResponsePinnedPostListType struct {
	Status         int                         `json:"status"`
	Code           string                      `json:"code"`
	Message        string                      `json:"message"`
	Result         bool                        `json:"result"`
	PinnedPostList []SelectAllPostDataResponse `json:"pinnedPostList"`
	PostCount      string                      `json:"postCount"`
	Page           int                         `json:"page"`
	Size           int                         `json:"size"`
}

// 게시글 리스트 응답 구조체
type ResponsePostByTagListType struct {
	Status    int                       `json:"status"`
	Code      string                    `json:"code"`
	Message   string                    `json:"message"`
	Result    bool                      `json:"result"`
	PostList  []PostsByTagsResponseType `json:"postList"`
	PostCount string                    `json:"postCount"`
}

// 카테고리로 게시글 조회
type ResponsePostByCategoryListType struct {
	Status    int                          `json:"status"`
	Code      string                       `json:"code"`
	Message   string                       `json:"message"`
	Result    bool                         `json:"result"`
	PostList  []PostByCategoryResponseType `json:"postList"`
	PostCount string                       `json:"postCount"`
}

// 게시글 리스트 응답 구조체
type ResponseInsertIdType struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Result   bool   `json:"result"`
	InsertId string `json:"insertId"`
}
