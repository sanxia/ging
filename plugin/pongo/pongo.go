package pongo

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin/render"
)

/* ================================================================================
 * Pongo模版引擎模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
type (
	RenderOptions struct {
		TemplatePath string
		Extensions   []string
		ContentType  string
	}

	PongoRender struct {
		Template *pongo2.Template
		Data     interface{}
		Options  *RenderOptions
	}
)

var (
	templateCache map[string]*pongo2.Template
	renderOptions RenderOptions
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Pongo模版渲染引擎
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRender(args ...string) *PongoRender {
	templatePath := "./templates"
	if len(args) > 0 {
		templatePath = args[0]
	}
	return New(RenderOptions{
		TemplatePath: templatePath,
		Extensions:   []string{".tpl", ".tmpl"},
		ContentType:  "text/html; charset=utf-8",
	})
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取自定义的Pongo模版渲染引擎
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func New(options RenderOptions) *PongoRender {
	setRenderOptions(options)

	compileTemplates()

	return &PongoRender{
		Options: &renderOptions,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取模版渲染器接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (p PongoRender) Instance(templateName string, data interface{}) render.Render {
	templateFilename, _ := filepath.Abs(path.Join(p.Options.TemplatePath, templateName))
	tmplKey := Md5(templateFilename)

	log.Printf("tmpl instance tmplKey: %s, tmplName: %s", tmplKey, templateFilename)

	template, exist := templateCache[tmplKey]
	if !exist {
		loadTemplate(templateFilename, *p.Options)
		template, exist = templateCache[tmplKey]
		if !exist {
			errors.New("not found template file: " + templateFilename)
		}
	} else {
		log.Printf("tmpl file: %s hits...success", templateFilename)
	}

	return PongoRender{
		Template: template,
		Data:     data,
		Options:  p.Options,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 渲染接口实现
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (p PongoRender) Render(w http.ResponseWriter) error {
	writeContentType(w, []string{p.Options.ContentType})
	err := p.Template.ExecuteWriter(p.Data.(map[string]interface{}), w)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置渲染选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func setRenderOptions(options RenderOptions) {
	if len(options.TemplatePath) == 0 {
		options.TemplatePath = "./templates"
	}
	if len(options.Extensions) == 0 {
		options.Extensions = []string{".tpl", ".tmpl"}
	}

	renderOptions = options
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 编译模版
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func compileTemplates() {
	templateCache = make(map[string]*pongo2.Template)

	allTemplateFiles := getTemplateFilenames(renderOptions.TemplatePath, renderOptions.Extensions)
	for _, fullfilename := range allTemplateFiles {
		filename := path.Base(fullfilename)
		extName := path.Ext(filename)
		tmplFilename := strings.TrimSuffix(fullfilename, extName)

		loadTemplate(tmplFilename, renderOptions)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 加载模版文件集合
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func getTemplateFilenames(templateRoot string, extNames []string) []string {
	templateRoot, _ = filepath.Abs(templateRoot)
	log.Printf("getTemplateFilenames templateRoot %s\n", templateRoot)

	filenames := make([]string, 0)
	err := filepath.Walk(templateRoot, func(fullFilePath string, fileInfo os.FileInfo, err error) error {
		filename := path.Base(fileInfo.Name())
		for _, ext := range extNames {
			if strings.HasSuffix(filename, ext) {
				filenames = append(filenames, fullFilePath)
				break
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("filepath.Walk() returned %v\n", err)
	}

	return filenames
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 加载模版文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func loadTemplate(templatePath string, options RenderOptions) {
	founded := true
	templateFilePath, templateFilename := templatePath, templatePath

	if !path.IsAbs(templatePath) {
		templateFilePath, _ = filepath.Abs(path.Join(options.TemplatePath, templatePath))
		templateFilename = templateFilePath
	}

	for _, extName := range options.Extensions {
		if !strings.HasSuffix(templatePath, extName) {
			templateFilePath = templateFilename + extName
		}

		if isMatched, _ := regexp.MatchString(".*"+extName+"$", templateFilePath); isMatched {
			if isExists := FileIsExists(templateFilePath); !isExists {
				continue
			}

			if tmpl, err := pongo2.FromFile(templateFilePath); err != nil {
				log.Fatalf("error:  \"%s\": %v", templateFilename, err)
			} else {
				tmplKey := Md5(templateFilename)
				log.Printf("load template: %s, %s", tmplKey, templateFilePath)
				templateCache[tmplKey] = tmpl
				founded = true
				break
			}
		} else {
			founded = false
		}
	}
	if !founded {
		errMsg := "template: " + templateFilePath + " not found"
		log.Fatalf("\"%s\"", errMsg)
	}
}
