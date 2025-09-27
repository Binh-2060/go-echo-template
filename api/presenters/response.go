package presenters

import (
	"time"
)

const (
	SUCCESS = 1
	FAIL    = 0
)

/*
Return success response - (can be modify as you want)
*/
func ResponseSuccess(data interface{}) map[string]interface{} {
	t := time.Now()
	return map[string]interface{}{
		"timestamp": t.Format("2006-01-02-15-04-05"),
		"status":    SUCCESS,
		"items":     data,
		"error":     nil,
	}
}

/*
Return list data with pagination infos
*/
// NOTE: if not pagination infos is not required, pass -1 to currentPage, currentPageTotalItem, totalPage
func ResponseSuccessListData(data interface{}, currentPage, currentPageTotalItem, totalPage int) map[string]interface{} {
	t := time.Now()
	return map[string]interface{}{
		"timestamp": t.Format("2006-01-02-15-04-05"),
		"status":    SUCCESS,
		"items": map[string]interface{}{
			"list_data": data,
			"pagination": map[string]interface{}{
				"current_page":            currentPage,
				"current_page_total_item": currentPageTotalItem,
				"total_page":              totalPage,
			},
		},
		"error": nil,
	}
}
