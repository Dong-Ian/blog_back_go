package utils

import "context"

func GetBlogIdFromContext(ctx context.Context) (string, bool) {
	blogId, ok := ctx.Value("BlogId").(string)
	return blogId, ok
}
