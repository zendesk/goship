package resources

import (
	"text/template"
)

// ParseTemplate creates template from template string
func ParseTemplate(templateContent string) *template.Template {
	return template.Must(
		template.New("tmpl").Funcs(
			template.FuncMap{
				"String": func(s *string) string {
					if s != nil {
						return *s
					}
					return ""
				},
			},
		).Parse(templateContent),
	)
}
