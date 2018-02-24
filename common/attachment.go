package common

/* ================================================================================
 * 附件数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type AttachmentList []*Attachment
type Attachment struct {
	Filename string `form:"filename" json:"filename"` //原始文件名（test.jpg）
	Type     string `form:"type" json:"type"`         //文件类型编码（text:文本 | image: 图片 | audio:音频 | video:视频 | other:其它）
	Size     int64  `form:"size" json:"size"`         //大小（单位：字节）
	Duration int64  `form:"duration" json:"duration"` //时长（单位：秒）
	Path     string `form:"path" json:"path"`         //全路径（本地磁盘或第三方文件系统）
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 上传附件选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type UploadAttachmentOption struct {
	FileSize  uint32   `form:"file_size" json:"file_size"`   //允许上传单个附件大小（单位：字节）
	FileCount uint32   `form:"file_count" json:"file_count"` //允许上传附件个数
	FileType  []string `form:"file_type" json:"file_type"`   //允许上传附件类型（rar,jpg,acc）
	IsEnabled bool     `form:"is_enabled" json:"is_enabled"` //是否允许上传
}
