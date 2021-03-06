package gopi

import (
	"html/template"
	//"log"
	//"unsafe"
	//"net/http"
	"time"
	"fmt"
	"bytes"
)

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

	"slice": func(renderArgs map[string]interface{}, key string, value ...interface{}) interface{} {
		renderArgs[key] = value
		return nil
	},

	"param": func(key string) template.HTML {
		if val, ok := app.Params[key]; ok {
			return template.HTML(template.HTMLEscapeString(val))
		} else {
			return ""
		}
	},

	"config": func(key string) bool {
		if val, ok := app.Params[key]; ok {
			if val == "yes" {
				return true
			} else {
				return false
			}
		}
		return false
	},

	"ListView": func(view string, list []interface{}) template.HTML {
		var html bytes.Buffer
		for index, item := range list {
			out, err := app.Render.execute(view, RenderArgs{"item": item, "index": index})
			if err != nil {
				return "LIST VIEW ERROR:\n" + template.HTML(err.Error()) + "\n\n"
			}

			html.Write(out.Bytes())
		}

		return template.HTML(html.String())
	},

	"Pager": func(pagination *Pagination, cssclass string) template.HTML {
		if cssclass != "" {
			pagination.css = cssclass
		}

		return pagination.Html()
	},

	"PagerSEO": func(pagination *Pagination) template.HTML {
		return pagination.SEO()
	},

	"raw": func(text string) template.HTML {
		return template.HTML(text)
	},

//	"option" : func(values []interface{})

	// Year
	"year": func() int {
		return time.Now().Year();
	},

	// Format a date according to the application's default date(time) format.
	"date": func(date time.Time) string {
		return date.Format(DateFormat)
	},
	"datetime": func(date time.Time) string {
		return date.Format(DateTimeFormat)
	},

	//

	"empty": func(text string) bool {
		return len(text) == 0
	},

	"not_empty": func(text string) bool {
		return len(text) > 0
	},

	// Forms

	"formError": func(field string, errorsMap map[string]string) template.HTML {
		if value, found := errorsMap[field]; found {
			return template.HTML("<div class=\"error\">"+value+"</div>")
		} else {
			return template.HTML("")
		}
	},

	"hasError": func(field string, errorsMap map[string]string) template.HTML {
		if _, found := errorsMap[field]; found {
			return template.HTML("has-error")
		} else {
			return template.HTML("")
		}
	},

	//
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called with no layout defined")
	},
}

// Included helper functions for use when rendering HTML.
var helperFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called with no layout defined")
	},
	"block": func() (string, error) {
		return "", fmt.Errorf("block called with no layout defined")
	},
	"current": func() (string, error) {
		return "", nil
	},
}