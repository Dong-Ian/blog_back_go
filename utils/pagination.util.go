package utils

import "strconv"

func ParsePaginationParams(pageSizeStr, pageNumberStr string) (offset, limit int) {
	const (
		defaultPageSize = 10
		defaultPageNum  = 1
	)

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = defaultPageSize
	}

	pageNum, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNum < 1 {
		pageNum = defaultPageNum
	}

	offset = (pageNum - 1) * pageSize
	limit = pageSize

	return offset, limit
}
