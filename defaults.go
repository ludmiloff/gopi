package gopi

const (
	DEFAULT_DATE_FORMAT     = "2006-01-02"
	DEFAULT_DATETIME_FORMAT = "2006-01-02 15:04"
)

var (
	DateFormat     string
	DateTimeFormat string
)

func (this *Application) SetDefaults() {
	DateFormat = DEFAULT_DATE_FORMAT
	DateTimeFormat = DEFAULT_DATETIME_FORMAT
}