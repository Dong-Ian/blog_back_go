package types

type UserGetProfileRequest struct {
	UserId string `json:"userId"`
	BlogId string `json:"blogId"`
}

// 프로필 색 변경 요청
 type UserUpdateProfileColorRequest struct {
	Color string `json:"color"`
 }

 // 블로그 타이틀 변경 요청
 type UserUpdateBlogTitleRequest struct {
	Title string `json:"title"`
 }

  // 프로필 수정
 type UserChangeProfileRequest struct {
	Name string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
	Title string `json:"title,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	GithubUrls string `json:"githubUrl,omitempty"`
	PersonalUrls string `json:"personalUrl,omitempty"`
	Memo string `json:"memo,omitempty"`
	// BackgroundImage string `json:"backgroundImage" binding:"optional"`
	// ProfileImage string `json:"profileImage" binding:"optional"`
 }
