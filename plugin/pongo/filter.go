package pongo

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

import (
	"github.com/dustin/go-humanize"
	"github.com/extemporalgenome/slug"
	"github.com/flosch/pongo2"
	"github.com/russross/blackfriday"
)

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Pongo模版引擎过滤器模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注册模版过滤器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func init() {
	//custom
	pongo2.RegisterFilter("sxcontains", filterContains)
	pongo2.RegisterFilter("sxcontainsany", filterContainsAny)
	pongo2.RegisterFilter("sxcontainsall", filterContainsAll)
	pongo2.RegisterFilter("sxhtmlencode", filterHtmlEncode)
	pongo2.RegisterFilter("sxhtml", filterHtmlDecode)
	pongo2.RegisterFilter("sxbr", filterLine2Br)
	pongo2.RegisterFilter("sxthumb", filterThumbImage)
	pongo2.RegisterFilter("sxdomain", filterUserDomain)

	// Regulars
	pongo2.RegisterFilter("slugify", filterSlugify)
	pongo2.RegisterFilter("filesizeformat", filterFilesizeformat)
	pongo2.RegisterFilter("truncatesentences", filterTruncatesentences)
	pongo2.RegisterFilter("truncatesentences_html", filterTruncatesentencesHtml)

	// Markup
	pongo2.RegisterFilter("markdown", filterMarkdown)

	// Humanize
	pongo2.RegisterFilter("naturalday", filterNaturalday)
	pongo2.RegisterFilter("intcomma", filterIntcomma)
	pongo2.RegisterFilter("ordinal", filterOrdinal)

	pongo2.RegisterFilter("local_time", ToLocalTime)
	pongo2.RegisterFilter("location_time", ToLocationTime)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * html 编码
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func filterHtmlEncode(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	newString := glib.HtmlEncode(in.String())
	return pongo2.AsSafeValue(newString), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * html 解码
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func filterHtmlDecode(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	newString := glib.HtmlDecode(in.String())
	return pongo2.AsSafeValue(newString), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * \r\n 和 \n 转 html br标签
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func filterLine2Br(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	newString := glib.String2Br(in.String())
	return pongo2.AsSafeValue(newString), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 是否包含一个
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func filterContains(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	inStringSlice := make([]string, 0)
	paramString := ""
	isContains := false

	if inInter := in.Interface(); inInter != nil {
		if inSlice, isOk := inInter.([]string); isOk {
			inStringSlice = inSlice
		}
	}

	if inputInter := param.Interface(); inputInter != nil {
		if inputValue, isOk := inputInter.(string); isOk {
			paramString = inputValue
		}
	}

	for _, inItemString := range inStringSlice {
		if strings.ToLower(inItemString) == strings.ToLower(paramString) {
			isContains = true
			break
		}
	}

	return pongo2.AsSafeValue(isContains), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 是否包含任意一个
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func filterContainsAny(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	inStringSlice := make([]string, 0)
	paramStringSlice := make([]string, 0)
	isContains := false

	if inInter := in.Interface(); inInter != nil {
		if inSlice, isOk := inInter.([]string); isOk {
			inStringSlice = inSlice
		}
	}

	if paramInter := param.Interface(); paramInter != nil {
		if paramString, isOk := paramInter.(string); isOk {
			paramStringSlice = strings.Split(paramString, ",")
		}
	}

	for _, paramItemString := range paramStringSlice {
		isFind := false

		paramItemString = strings.ToLower(paramItemString)
		for _, inItemString := range inStringSlice {
			if paramItemString == strings.ToLower(inItemString) {
				isFind = true
				break
			}
		}

		if isFind {
			isContains = true
			break
		}
	}

	return pongo2.AsSafeValue(isContains), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 是否全部都包含
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func filterContainsAll(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	inStringSlice := make([]string, 0)
	paramStringSlice := make([]string, 0)
	isContains := true

	if inInter := in.Interface(); inInter != nil {
		if inSlice, isOk := inInter.([]string); isOk {
			inStringSlice = inSlice
		}
	}

	if paramInter := param.Interface(); paramInter != nil {
		if paramString, isOk := paramInter.(string); isOk {
			paramStringSlice = strings.Split(paramString, ",")
		}
	}

	for _, paramItemString := range paramStringSlice {
		isFind := false

		paramItemString = strings.ToLower(paramItemString)
		for _, inItemString := range inStringSlice {
			if paramItemString == strings.ToLower(inItemString) {
				isFind = true
				break
			}
		}

		if !isFind {
			isContains = false
			break
		}
	}

	return pongo2.AsSafeValue(isContains), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 缩略图地址处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func filterThumbImage(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	sourceString := in.String()
	if sourceString == "" || !strings.Contains(sourceString, "img.test.com") {
		return pongo2.AsSafeValue(sourceString), nil
	}

	//大小映射
	sizeMap := map[string]string{
		"5050":   "50x50",
		"75100":  "75x100",
		"100100": "100x100",
		"150150": "150x150",
		"200200": "200x200",
		"300400": "300x400",
	}

	size := param.String()
	sizeValue := "50x50"
	if value, ok := sizeMap[size]; ok {
		sizeValue = value
	}

	newString := sourceString + "_" + sizeValue + ".jpg"

	return pongo2.AsSafeValue(newString), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 用户域名处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func filterUserDomain(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	username := in.String()
	userId := param.String()

	if username == "" {
		username = userId
	}

	domain := "http://a.test.com/" + username

	return pongo2.AsSafeValue(domain), nil
}

func filterMarkdown(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsSafeValue(string(blackfriday.MarkdownCommon([]byte(in.String())))), nil
}

func filterSlugify(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(slug.Slug(in.String())), nil
}

func filterFilesizeformat(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(humanize.IBytes(uint64(in.Integer()))), nil
}

var filterTruncatesentencesRe = regexp.MustCompile(`(?U:.*[\w]{3,}.*([\d][\.!?][\D]|[\D][\.!?][\s]|[\n$]))`)

func filterTruncatesentences(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	count := param.Integer()
	if count <= 0 {
		return pongo2.AsValue(""), nil
	}
	sentencens := filterTruncatesentencesRe.FindAllString(strings.TrimSpace(in.String()), -1)
	index := int(math.Min(float64(count), float64(len(sentencens))))
	return pongo2.AsValue(strings.TrimSpace(strings.Join(sentencens[:index], ""))), nil
}

// Taken from pongo2/filters_builtin.go
func filterTruncateHtmlHelper(value string, new_output *bytes.Buffer, cond func() bool, fn func(c rune, s int, idx int) int, finalize func()) {
	vLen := len(value)
	tag_stack := make([]string, 0)
	idx := 0
	for idx < vLen && !cond() {
		c, s := utf8.DecodeRuneInString(value[idx:])
		if c == utf8.RuneError {
			idx += s
			continue
		}
		if c == '<' {
			new_output.WriteRune(c)
			idx += s // consume "<"
			if idx+1 < vLen {
				if value[idx] == '/' {
					// Close tag
					new_output.WriteString("/")
					tag := ""
					idx += 1 // consume "/"
					for idx < vLen {
						c2, size2 := utf8.DecodeRuneInString(value[idx:])
						if c2 == utf8.RuneError {
							idx += size2
							continue
						}
						// End of tag found
						if c2 == '>' {
							idx++ // consume ">"
							break
						}
						tag += string(c2)
						idx += size2
					}
					if len(tag_stack) > 0 {
						// Ideally, the close tag is TOP of tag stack
						// In malformed HTML, it must not be, so iterate through the stack and remove the tag
						for i := len(tag_stack) - 1; i >= 0; i-- {
							if tag_stack[i] == tag {
								// Found the tag
								tag_stack[i] = tag_stack[len(tag_stack)-1]
								tag_stack = tag_stack[:len(tag_stack)-1]
								break
							}
						}
					}
					new_output.WriteString(tag)
					new_output.WriteString(">")
				} else {
					// Open tag
					tag := ""
					params := false
					for idx < vLen {
						c2, size2 := utf8.DecodeRuneInString(value[idx:])
						if c2 == utf8.RuneError {
							idx += size2
							continue
						}
						new_output.WriteRune(c2)
						// End of tag found
						if c2 == '>' {
							idx++ // consume ">"
							break
						}
						if !params {
							if c2 == ' ' {
								params = true
							} else {
								tag += string(c2)
							}
						}
						idx += size2
					}
					// Add tag to stack
					tag_stack = append(tag_stack, tag)
				}
			}
		} else {
			idx = fn(c, s, idx)
		}
	}
	finalize()
	for i := len(tag_stack) - 1; i >= 0; i-- {
		tag := tag_stack[i]
		// Close everything from the regular tag stack
		new_output.WriteString(fmt.Sprintf("</%s>", tag))
	}
}
func filterTruncatesentencesHtml(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	count := param.Integer()
	if count <= 0 {
		return pongo2.AsValue(""), nil
	}
	value := in.String()
	newLen := int(math.Max(float64(param.Integer()), 0))
	new_output := bytes.NewBuffer(nil)
	sentencefilter := 0
	filterTruncateHtmlHelper(value, new_output, func() bool {
		return sentencefilter >= newLen
	}, func(_ rune, _ int, idx int) int {
		// Get next word
		word_found := false
		for idx < len(value) {
			c2, size2 := utf8.DecodeRuneInString(value[idx:])
			if c2 == utf8.RuneError {
				idx += size2
				continue
			}
			if c2 == '<' {
				// HTML tag start, don't consume it
				return idx
			}
			new_output.WriteRune(c2)
			idx += size2
			if (c2 == '.' && !(idx+1 < len(value) && value[idx+1] >= '0' && value[idx+1] <= '9')) ||
				c2 == '!' || c2 == '?' || c2 == '\n' {
				// Sentence ends here, stop capturing it now
				break
			} else {
				word_found = true
			}
		}
		if word_found {
			sentencefilter++
		}
		return idx
	}, func() {})
	return pongo2.AsSafeValue(new_output.String()), nil
}

func filterIntcomma(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(humanize.Comma(int64(in.Integer()))), nil
}
func filterOrdinal(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(humanize.Ordinal(in.Integer())), nil
}

func filterNaturalday(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	basetime, is_time := in.Interface().(time.Time)
	if !is_time {
		return nil, &pongo2.Error{
			Sender:    "filter:naturalday",
			OrigError: errors.New("naturalday-value is not a time.Time-instance."),
		}
	}
	var reference_time time.Time
	if !param.IsNil() {
		reference_time, is_time = param.Interface().(time.Time)
		if !is_time {
			return nil, &pongo2.Error{
				Sender:    "filter:naturalday",
				OrigError: errors.New("naturalday-parameter is not a time.Time-instance."),
			}
		}
	} else {
		reference_time = time.Now()
	}
	d := reference_time.Sub(basetime) / time.Hour
	switch {
	case d >= 0 && d < 24:
		// Today
		return pongo2.AsValue("today"), nil
	case d >= 24:
		return pongo2.AsValue("yesterday"), nil
	case d < 0 && d >= -24:
		return pongo2.AsValue("tomorrow"), nil
	}
	// Default behaviour
	return pongo2.ApplyFilter("naturaltime", in, param)
}

func ToLocalTime(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, isTime := in.Interface().(time.Time)
	if !isTime {
		return nil, &pongo2.Error{
			Sender:    "filter:time",
			OrigError: errors.New("Filter argument must be of type 'time.Time'."),
		}
	}
	loc, _ := time.LoadLocation("Local")
	return pongo2.AsValue(t.In(loc).Format("15:04")), nil
}

func ToLocationTime(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, isTime := in.Interface().(time.Time)
	if !isTime {
		return nil, &pongo2.Error{
			Sender:    "filter:localtime",
			OrigError: errors.New("Filter input argument must be of type 'time.Time'."),
		}
	}

	loc, _ := time.LoadLocation(param.String())
	return pongo2.AsValue(t.In(loc).Format("15:04")), nil
}
