package gopi

import (
	"bytes"
	//"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	//"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ContentType header constant.
	ContentType = "Content-Type"
	// ContentLength header constant.
	ContentLength = "Content-Length"
	// ContentBinary header value for binary data.
	ContentBinary = "application/octet-stream"
	// ContentJSON header value for JSON data.
	ContentJSON = "application/json"
	// ContentHTML header value for HTML data.
	ContentHTML = "text/html"
	// ContentXHTML header value for XHTML data.
	ContentXHTML = "application/xhtml+xml"
	// ContentXML header value for XML data.
	ContentXML = "text/xml"
	// Default character encoding.
	defaultCharset = "UTF-8"
)

type RenderArgs map[string]interface{}

// Included helper functions for use when rendering html.
var TemplateFunctions = template.FuncMap{
	"set": func(renderArgs map[string]interface{}, key string, value interface{}) interface{} {
		renderArgs[key] = value
		return nil
	},

	"append": func(renderArgs map[string]interface{}, key string, value interface{}) interface{} {
		if renderArgs[key] == nil {
			renderArgs[key] = []interface{}{value}
		} else {
			renderArgs[key] = append(renderArgs[key].([]interface{}), value)
		}
		return nil
	},

	"param": func(key string) template.HTML {
		if val, ok := App.Params[key]; ok {
			return template.HTML(template.HTMLEscapeString(val))
		} else {
			return ""
		}
	},
}

// Delims represents a set of Left and Right delimiters for HTML template rendering.
type RenderDelims struct {
	// Left delimiter, defaults to {{.
	Left string
	// Right delimiter, defaults to }}.
	Right string
}

// Options is a struct for specifying configuration options for the render.Render object.
type RenderOptions struct {
	// Directory to load templates. Default is "templates".
	Directory string
	// Layout template name. Will not render a layout if blank (""). Defaults to blank ("").
	Layout string
	// Extensions to parse template files from. Defaults to [".tmpl"].
	Extensions []string
	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	Funcs template.FuncMap
	// Delims sets the action delimiters to the specified strings in the Delims struct.
	Delims RenderDelims
	// Appends the given character set to the Content-Type header. Default is "UTF-8".
	Charset string
	// Outputs human readable JSON.
	IndentJSON bool
	// Outputs human readable XML.
	IndentXML bool
	// Prefixes the JSON output with the given bytes.
	PrefixJSON []byte
	// Prefixes the XML output with the given bytes.
	PrefixXML []byte
	// Allows changing of output to XHTML instead of HTML. Default is "text/html"
	HTMLContentType string
	// If IsDevelopment is set to true, this will recompile the templates on every request. Default if false.
	IsDevelopment bool
}

// HTMLOptions is a struct for overriding some rendering Options for specific HTML call.
type HTMLOptions struct {
	// Layout template name. Overrides Options.Layout.
	Layout string
}

// Head defines the basic ContentType and Status fields.
type Head struct {
	ContentType string
	Status      int
}

// XML built-in renderer.
type XML struct {
	Head
	Indent bool
	Prefix []byte
}

// JSON built-in renderer.
type JSON struct {
	Head
	Indent bool
	Prefix []byte
}

const (
	POS_HEAD  = 1
	POS_END   = 2
	POS_READY = 3
)

type ScriptArray []string

type HTML struct {
	PageTitle string
	Meta      map[string]string
	Css       map[uint]string
	Scripts   map[uint]ScriptArray
	//
	css uint
}

// Data built-in renderer.
type Data struct {
	Head
}

// Render is a service that provides functions for easily writing JSON, XML,
// Binary Data, and HTML templates out to a http Response.
type Render struct {
	// Customize Secure with an Options struct.
	opt             RenderOptions
	templates       *template.Template
	compiledCharset string
}

// Initialize HTML struct
func NewHTML(title string) *HTML {
	h := &HTML{PageTitle: title}
	return h
}

// Write outputs the header content.
func (h Head) Write(w http.ResponseWriter) {
	w.Header().Set(ContentType, h.ContentType)
	w.WriteHeader(h.Status)
}

// New constructs a new Render instance with the supplied options.
func NewRender(options RenderOptions) *Render {
	r := Render{
		opt: options,
	}

	r.prepareOptions()
	r.compileTemplates()

	return &r
}

func (r *Render) prepareOptions() {
	// Fill in the defaults if need be.
	if len(r.opt.Charset) == 0 {
		r.opt.Charset = defaultCharset
	}
	r.compiledCharset = "; charset=" + r.opt.Charset

	if len(r.opt.Directory) == 0 {
		r.opt.Directory = "templates"
	}
	if len(r.opt.Extensions) == 0 {
		r.opt.Extensions = []string{".tmpl"}
	}
	if len(r.opt.HTMLContentType) == 0 {
		r.opt.HTMLContentType = ContentHTML
	}
}

func (r *Render) compileTemplates() {
	dir := r.opt.Directory
	r.templates = template.New(dir)
	r.templates.Delims(r.opt.Delims.Left, r.opt.Delims.Right)

	// Walk the supplied directory and compile any files that match our extension list.
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := ""
		if strings.Index(rel, ".") != -1 {
			ext = "." + strings.Join(strings.Split(rel, ".")[1:], ".")
		}

		//log.Println("TEMPLATE ", path)

		for _, extension := range r.opt.Extensions {
			if ext == extension {

				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := (rel[0 : len(rel)-len(ext)])
				tmpl := r.templates.New(filepath.ToSlash(name))

				// Add our funcmaps.
				// for _, funcs := range r.opt.Funcs {
				// 	tmpl.Funcs(funcs)
				// }

				tmpl.Funcs(r.opt.Funcs)

				text := string(buf)
				text = strings.Replace(text, "}}\n\n", "}}", -1)
				text = strings.Replace(text, "}}\n", "}}", -1)

				// Break out if this parsing fails. We don't want any silent server starts.
				template.Must(tmpl.Funcs(TemplateFunctions).Parse(text))
				break
			}
		}

		return nil
	})
}

// XML marshals the given interface object and writes the XML response.
// func (r *Render) XML(status int, v interface{}) {
// 	head := Head{
// 		ContentType: ContentXML + r.compiledCharset,
// 		Status:      status,
// 	}

// 	x := XML{
// 		Head:   head,
// 		Indent: r.opt.IndentXML,
// 		Prefix: r.opt.PrefixXML,
// 	}

// 	r.Render(x, v)
// }

// JSON marshals the given interface object and writes the JSON response.
// func (r *Render) JSON(status int, v interface{}) {
// 	head := Head{
// 		ContentType: ContentJSON + r.compiledCharset,
// 		Status:      status,
// 	}

// 	j := JSON{
// 		Head:   head,
// 		Indent: r.opt.IndentJSON,
// 		Prefix: r.opt.PrefixJSON,
// 	}

// 	r.Render(j, v)
// }

// Data writes out the raw bytes as binary data.
// func (r *Render) Data(w http.ResponseWriter, status int, v []byte) {
// 	head := Head{
// 		ContentType: ContentBinary,
// 		Status:      status,
// 	}

// 	d := Data{
// 		Head: head,
// 	}

// 	r.Render(d, v)
// }

// HTML builds up the response from the specified template and bindings.
func (r *Render) RenderHTML(w http.ResponseWriter, status int, name string, binding interface{}, htmlOpt ...HTMLOptions) {
	// If we are in development mode, recompile the templates on every HTML request.
	if r.opt.IsDevelopment {
		r.compileTemplates()
	}

	opt := r.prepareHTMLOptions(htmlOpt)

	// Assign a layout if there is one.
	if len(opt.Layout) > 0 {
		log.Println("USE YIELD")
		r.addYield(name, binding)
		name = opt.Layout
	}

	head := Head{
		ContentType: r.opt.HTMLContentType + r.compiledCharset,
		Status:      status,
	}

	out := new(bytes.Buffer)
	err := r.templates.ExecuteTemplate(out, name, binding)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	head.Write(w)
	w.Write(out.Bytes())
}

func (r *Render) execute(name string, binding interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	return buf, r.templates.ExecuteTemplate(buf, name, binding)
}

func (r *Render) addYield(name string, binding interface{}) {
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf, err := r.execute(name, binding)
			// Return safe HTML here since we are rendering our own template.
			return template.HTML(buf.String()), err
		},
		"current": func() (string, error) {
			return name, nil
		},

		// "set": func(renderArgs map[string]interface{}, key string, value interface{}) template.JS {
		// 	renderArgs[key] = value
		// 	return template.JS("")
		// },
	}
	r.templates.Funcs(funcs)
}

func (r *Render) prepareHTMLOptions(htmlOpt []HTMLOptions) HTMLOptions {
	if len(htmlOpt) > 0 {
		return htmlOpt[0]
	}

	return HTMLOptions{
		Layout: r.opt.Layout,
	}
}
