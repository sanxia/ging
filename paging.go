package ging

import (
	"math"
)

/* ================================================================================
 * paging
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	Paging struct {
		PagingIndex int64  `form:"paging_index" json:"paging_index"` //current page index
		PagingSize  int64  `form:"paging_size" json:"paging_size"`   //size per page
		TotalRecord int64  `form:"total_record" json:"total_record"` //total records
		PagingCount int64  `form:"paging_count" json:"paging_count"` //total pages
		Sortorder   string `form:"sortorder" json:"-"`               //sort
		Group       string `form:"group" json:"-"`                   //group
		IsTotalOnce bool   `form:"-" json:"-"`                       //calculate the total number of records only on the first page（every time by default）
	}
)

func NewPaging() *Paging {
	paging := new(Paging)
	paging.PagingIndex = 1
	paging.PagingSize = 10
	paging.PagingCount = 1
	paging.Sortorder = "id DESC"

	return paging
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set the total number of records
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (paging *Paging) SetTotalRecord(totalRecord int64) {
	paging.PagingCount = 1
	if totalRecord > 0 {
		paging.TotalRecord = totalRecord
		paging.PagingCount = int64(math.Ceil(float64(paging.TotalRecord) / float64(paging.PagingSize)))
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get paginated offset
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (paging *Paging) Offset() int64 {
	if paging.PagingIndex <= 1 || paging.PagingSize == 0 {
		return 0
	}

	offsetIndex := (paging.PagingIndex - 1) * paging.PagingSize

	if offsetIndex > paging.TotalRecord {
		offsetIndex = paging.TotalRecord
	}

	return offsetIndex
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get the end record index
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (paging *Paging) EndIndex() int64 {
	endIndex := paging.PagingSize
	if paging.PagingIndex > 1 {
		endIndex = paging.PagingIndex * paging.PagingSize
	}

	if endIndex > paging.TotalRecord {
		endIndex = paging.TotalRecord
	}

	return endIndex
}
