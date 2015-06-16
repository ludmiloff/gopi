package gopi

import (
	"html/template"
	//"log"
	//"unsafe"
	//"net/http"
	"time"
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

	"param": func(key string) template.HTML {
		if val, ok := App.Params[key]; ok {
			return template.HTML(template.HTMLEscapeString(val))
		} else {
			return ""
		}
	},

	"ListView": func(view string, list []interface{}) template.HTML {
		var html = ""
		for index, item := range list {
			out, err := App.Render.execute(view, RenderArgs{"item": item, "index": index})
			if err != nil {
				return "LIST VIEW ERROR:\n" + template.HTML(err.Error()) + "\n\n"
			}

			html = html + out.String()
		}

		return template.HTML(html)
	},

	"raw": func(text string) template.HTML {
		return template.HTML(text)
	},

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
}