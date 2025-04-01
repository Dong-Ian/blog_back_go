package types

// 게시글 삭제 요청
type DeletePostRequest struct {
	PostSeq string `json:"postSeq"`
}
