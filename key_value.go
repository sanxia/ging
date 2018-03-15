package ging

/* ================================================================================
 * 键值对数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type KeyValuePairList []*KeyValuePair
type KeyValuePair struct {
	Key   string `form:"key" json:"key"`
	Value string `form:"value" json:"value"`
}

type KeyValueNodeList []*KeyValueNode
type KeyValueNode struct {
	KeyValuePair
	Childs KeyValueNodeList `form:"childs" json:"childs"`
}
