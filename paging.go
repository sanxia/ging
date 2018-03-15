package ging

import (
	"math"
)

/* ================================================================================
 * 分页数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	Paging struct {
		PagingIndex int64  `form:"paging_index" json:"paging_index"` //当前页索引
		PagingSize  int64  `form:"paging_size" json:"paging_size"`   //每页大小
		TotalRecord int64  `form:"total_record" json:"total_record"` //总记录数
		PagingCount int64  `form:"paging_count" json:"paging_count"` //总页数
		Sortorder   string `form:"sortorder" json:"-"`               //排序
		Group       string `form:"group" json:"-"`                   //分组
		IsTotalOnce bool   `form:"-" json:"-"`                       //是否只在第一页计算总记录数（默认每次都计算）
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
 * 设置总记录数
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (paging *Paging) SetTotalRecord(totalRecord int64) {
	paging.PagingCount = 1
	if totalRecord > 0 {
		paging.TotalRecord = totalRecord
		paging.PagingCount = int64(math.Ceil(float64(paging.TotalRecord) / float64(paging.PagingSize)))
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取分页偏移
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (paging *Paging) Offset() int64 {
	if paging.PagingIndex <= 1 || paging.PagingSize == 0 {
		return 0
	}

	offset := (paging.PagingIndex - 1) * paging.PagingSize
	return offset
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取结束记录索引
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (paging *Paging) EndIndex() int64 {
	if paging.PagingIndex <= 1 {
		return paging.PagingSize
	}

	offset := paging.PagingIndex * paging.PagingSize
	return offset
}
