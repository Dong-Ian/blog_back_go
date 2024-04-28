package types

// 상세 포스트 조회 요청
type ViewPostContents struct {
	PostSeq string `json:"postSeq" binding:"required"`
}

// 특정 게시글 응답
type ViewSpecificPostContentsResponse struct {
	PostSeq string `json:"postSeq"`
	PostTitle string `json:"postTitle"`
	PostContents string `json:"postContents"`
	Tags []string `json:"tags"`
	UserId string `json:"userId"`
	UserName string `json:"userName"`
	Viewed string `json:"viewed"`
	IsPinned string `json:"isPinned"`
	Urls []string `json:"urls"`
	RegDate string `json:"regDate"`
	ModDate string `json:"modDate"`
}

// 특정 게시글 태그 데이터 쿼리 결과
type SelectSpeicificPostTagDataResult struct {
	TagName string
}

// 특정 게시글 데이터 쿼리 결과
type SelectSpecificPostDataResult struct {
	PostSeq string
	PostTitle string
	PostContents string
	PostStatus string
	Tags string
	CategoryName string
	UserId string
	UserName string
	Viewed string
	IsPinned string
	RegDate string
	ModDate string
}


// 태그로 특정 게시글 쿼리
type SelectPostsByTags struct {
	TagName string		
	CategoryName string
	UserName string
	PostSeq string		
	PostTitle string	
	PostContents string
	Viewed string	
	RegDate string	
	ModDate string 
}

