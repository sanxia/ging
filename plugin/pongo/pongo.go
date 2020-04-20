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
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Pongo模版引擎模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	PongoTemplate struct {
		Template *pongo2.Template
		Data     interface{}
		Option   *PongoOption
	}

	PongoOption struct {
		TemplatePath string
		Extensions   []string
		ContentType  string
	}
)

var (
	templateCache map[string]*pongo2.Template
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * initializer pongo template engine
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewPongoTemplate(args ...*PongoOption) *PongoTemplate {
	var option *PongoOption

	if len(args) > 0 {
		option = args[0]
	}

	if option == nil {
		option = new(PongoOption)
	}

	if len(option.TemplatePath) == 0 {
		option.TemplatePath = "./templates"
	}

	if len(option.Extensions) == 0 {
		option.Extensions = []string{".tpl", ".tmpl"}
	}

	if len(option.ContentType) == 0 {
		option.ContentType = "text/html; charset=utf-8"
	}

	pongoTemplate := &PongoTemplate{
		Option: option,
	}

	pongoTemplate.CompileTemplate()

	return pongoTemplate
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gin render interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (p PongoTemplate) Instance(templateName string, data interface{}) render.Render {
	templateFilename, _ := filepath.Abs(path.Join(p.Option.TemplatePath, templateName))
	tmplKey := glib.Md5(templateFilename)

	template, exist := templateCache[tmplKey]
	if !exist {
		p.LoadTemplate(templateFilename)
		template, exist = templateCache[tmplKey]
		if !exist {
			errors.New("not found template file: " + templateFilename)
		}
	} else {
		log.Printf("tmpl file: %s hits...success", templateFilename)
	}

	return PongoTemplate{
		Template: template,
		Data:     data,
		Option:   p.Option,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Render interface impl
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (p PongoTemplate) Render(w http.ResponseWriter) error {
	p.WriteContentType(w)

	err := p.Template.ExecuteWriter(p.Data.(map[string]interface{}), w)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Render interface impl
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (p PongoTemplate) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{p.Option.ContentType}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * compile template file
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (p PongoTemplate) CompileTemplate() {
	templateCache = make(map[string]*pongo2.Template)

	templateFilenames := getTemplateFilenames(p.Option.TemplatePath, p.Option.Extensions)

	for _, fullfilename := range templateFilenames {
		filename := path.Base(fullfilename)
		extName := path.Ext(filename)
		tmplFilename := strings.TrimSuffix(fullfilename, extName)

		p.LoadTemplate(tmplFilename)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * load template file
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (p PongoTemplate) LoadTemplate(templatePath string) {
	founded := true
	templateFilePath, templateFilename := templatePath, templatePath

	if !path.IsAbs(templatePath) {
		templateFilePath, _ = filepath.Abs(path.Join(p.Option.TemplatePath, templatePath))
		templateFilename = templateFilePath
	}

	for _, extName := range p.Option.Extensions {
		if !strings.HasSuffix(templatePath, extName) {
			templateFilePath = templateFilename + extName
		}

		if isMatched, _ := regexp.MatchString(".*"+extName+"$", templateFilePath); isMatched {
			if isExists := glib.FileIsExists(templateFilePath); !isExists {
				continue
			}

			if tmpl, err := pongo2.FromFile(templateFilePath); err != nil {
				log.Fatalf("error:  \"%s\": %v", templateFilename, err)
			} else {
				tmplKey := glib.Md5(templateFilename)
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

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * render template string
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func RenderString(templateString string, data interface{}) (string, error) {
	tpl, err := pongo2.FromString(templateString)
	if err != nil {
		panic(err)
	}

	result, err := tpl.Execute(data.(map[string]interface{}))
	if err != nil {
		panic(err)
	}

	return result, err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * load template disk files
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
