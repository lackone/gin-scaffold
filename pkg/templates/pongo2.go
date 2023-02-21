package templates

import (
	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
	"path"
)

type Pongo2Ctx = pongo2.Context

type Pongo2Option struct {
	TemplateDir string
	TemplateSet *pongo2.TemplateSet
	ContentType string
}

type Pongo2Render struct {
	Option   *Pongo2Option
	Template *pongo2.Template
	Context  pongo2.Context
}

func NewPongo2Render(option *Pongo2Option) *Pongo2Render {
	return &Pongo2Render{
		Option: option,
	}
}

func NewDefaultPongo2Render() *Pongo2Render {
	return NewPongo2Render(&Pongo2Option{
		TemplateDir: "resources/views",
		TemplateSet: nil,
		ContentType: "text/html; charset=utf-8",
	})
}

func (p Pongo2Render) Instance(name string, data any) render.Render {
	var template *pongo2.Template

	filename := path.Join(p.Option.TemplateDir, name)

	if p.Option.TemplateSet != nil {
		template = pongo2.Must(p.Option.TemplateSet.FromFile(filename))
		data.(pongo2.Context).Update(p.Option.TemplateSet.Globals)
	} else {
		if gin.Mode() == "debug" {
			template = pongo2.Must(pongo2.FromFile(filename))
		} else {
			template = pongo2.Must(pongo2.FromCache(filename))
		}
	}

	return Pongo2Render{
		Template: template,
		Context:  data.(pongo2.Context),
		Option:   p.Option,
	}
}

func (p Pongo2Render) Render(w http.ResponseWriter) error {
	p.WriteContentType(w)
	err := p.Template.ExecuteWriter(p.Context, w)
	return err
}

func (p Pongo2Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{p.Option.ContentType}
	}
}
