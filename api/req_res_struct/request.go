package req_res

type PostCreateClassroom struct {
	Name        string `form:"name" json:"name"`
	CoverImage  string `form:"coverImage" json:"coverImageUrl"`
	Code        string `form:"code" json:"code" `
	Description string `form:"description" json:"description" `
}
