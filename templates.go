package gopi

import (
	"github.com/pelletier/go-toml"
	"html/template"
	"log"
)

func (this *Application) InitRender() {
	general := this.Config.Get("general").(*toml.TomlTree)
	log.Println("TEMPLATE PATH", general.Get("template_path"))
	this.Render = NewRender(RenderOptions{
		Directory:       general.GetDefault("template_path", "templates").(string),
		Layout:          "layout",
		Extensions:      []string{".tmpl", ".html"},
		Funcs:           []template.FuncMap{},
		Delims:          RenderDelims{"{{", "}}"},
		Charset:         "UTF-8",
		IndentJSON:      true,
		IndentXML:       true,
		PrefixJSON:      []byte(")]}',\n"),
		PrefixXML:       []byte("<?xml version='1.0' encoding='UTF-8'?>"),
		HTMLContentType: ContentHTML,
		IsDevelopment:   true,
	})
}
