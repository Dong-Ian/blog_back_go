package types

type EditPostRequest struct {
	PostSeq string `json:"postSeq"`
	PostTitle string `json:"postTitle"`
	PostContents string `json:"postContents"`
	Category 	string `json:"category,omitempty"`
	ImageSeqs  []string `json:"imageSeqs,omitempty"`
	Tags []string `json:"tags,omitempty"`
	IsPinned string `json:"isPinned"`
	// IsSecret string `json:"isSecret"`
}
