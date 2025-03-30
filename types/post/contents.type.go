package types

// 특정 게시글 응답
type ViewSpecificPostContentsResponse struct {
	PostSeq      string   `json:"postSeq"`
	PostTitle    string   `json:"postTitle"`
	PostContents string   `json:"postContents"`
	Tags         []string `json:"tags"`
	CategoryName string   `json:"category"`
	UserName     string   `json:"userName"`
	Viewed       string   `json:"viewed"`
	IsPinned     string   `json:"isPinned"`
	Urls         []string `json:"urls"`
	RegDate      string   `json:"regDate"`
	ModDate      string   `json:"modDate"`
}

type ViewImageUrl struct {
	Urls []string `json:"urls"`
}

// 특정 게시글 태그 데이터 쿼리 결과
// type SelectSpeicificPostTagDataResult struct {
// 	TagName *string
// }

// 특정 게시글 데이터 쿼리 결과
type SelectSpecificPostDataResult struct {
	PostSeq      string
	PostTitle    string
	PostContents string
	PostStatus   string
	Tags         *string
	CategoryName *string
	UserName     string
	Viewed       string
	IsPinned     string
	RegDate      string
	ModDate      string
}

// 태그로 특정 게시글 쿼리
type SelectPostsByTags struct {
	TagName      string
	CategoryName string
	UserName     string
	PostSeq      string
	PostTitle    string
	PostContents string
	Viewed       string
	RegDate      string
	ModDate      string
}
