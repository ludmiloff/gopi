package gopi
import (
	"html/template"
	"net/url"
	//"log"
	"strconv"
	"math"
	//"fmt"
	"bytes"
)

const (
	PAGER_MAX_PAGES = 10
)

type Pagination struct {
	current				int			// current page
	total_items			int			// total number of items
	num_pages			int 		// overall number of pages
	per_page			int			// max number of items on single page
	max					int			// widget maximum pages
	css					string		// widget css class
	url					string
	first_url			string
	First				string
	Last				string
	Prev				string
	Next				string
	Page_var			string		// page variable, used in url
	next_to_current		string
	prev_to_current		string
	prefix				string
}

func NewDefaultPagination(total_items, per_page int, current int, url *url.URL, prefix string) *Pagination {
	max := PAGER_MAX_PAGES // default
	n := int(math.Ceil(float64(total_items) / float64(per_page)))
	//n = int(math.Min(float64(n), float64(max)))
	//println("num pages = ", n)
	p := Pagination{
		current: 		current + 1,
		total_items:	total_items,
		num_pages: 		n,
		per_page:		per_page,
		max: 			max,
		css: 			"pagination",
		First: 			"",
		Last: 			"",
		Prev: 			"<<",
		Next: 			">>",
		Page_var: 		"page",
		prefix:			prefix,
	}
	if total_items > per_page {
		p.excludeParamFromUrl(url, "page")
	}
	return &p
}

func (this *Pagination) Page() int {
	if this.current > this.num_pages {
		this.current = this.num_pages
	} else  if this.current <= 0 {
		this.current = 1
	}

	return this.current
}

func (this *Pagination) Next2Current() string {
	if this.next_to_current != "" {
		return this.next_to_current
	}
	this.next_to_current = strconv.Itoa(this.current + 1)
	return this.next_to_current
}

func (this *Pagination) Prev2Current() string {
	if this.prev_to_current != "" {
		return this.prev_to_current
	}
	this.prev_to_current = strconv.Itoa(this.current - 1)
	return this.prev_to_current
}


func (this *Pagination) excludeParamFromUrl(url *url.URL, page string) {
	query := url.Query()
	query.Del(page)
	encoded := query.Encode()
	if encoded != "" {
		this.first_url = this.prefix + url.Path + "?" + encoded
		this.url = this.first_url + "&" + this.Page_var + "="
	} else {
		this.first_url = this.prefix + url.Path
		this.url = this.first_url + "?" + this.Page_var + "="
	}
}

func (this *Pagination) Pages() []int {
	if this.total_items > 0 {

		var pages []int

		max1 := this.max - 1
		page_nums := this.num_pages
		page := this.Page()
		//println("\npage_nums", page_nums)
		//println("page", page)
		mid := page_nums - 5
		switch  {
		case page > mid && page_nums > max1:
			//println("case 1")
			start := mid
			//println("start = ", start)
			pages = make([]int, 6)
			for i:= range pages {
				pages[i] = start + i
			}

		case page > 5 && page_nums > max1:
			//println("case 2")
			start := page - 4//5 + 1
			//println(start)
			n := int(math.Min(float64(max1), float64(page+5)))
			//println(n)
			//println(page+4+1)
			pages = make([]int, n)
			for i := range pages {
				pages[i] = start + i
			}

		default:
			pages = make([]int, int(math.Min(float64(max1), float64(page_nums))))
			for i := range pages {
				pages[i] = i + 1
			}
		}

		return pages
	} else {
		return nil
	}
}

func (this *Pagination) SEO() template.HTML {
	if this.total_items < this.per_page {
		return template.HTML("")
	}
	var html string = ""
	page := this.Page() // page = current
	if page == 2 {
		html = html + "<link rel=\"prev\" href=\"" + this.first_url +  "\" />\n"

	} else if page > 2 {
		url1 := this.url + this.Prev2Current() //strconv.Itoa(page - 1)
		html = html + "<link rel=\"prev\" href=\"" + url1 +  "\" />\n"
	}

	if this.current < this.num_pages {
		url1 := this.url + this.Next2Current() //strconv.Itoa(page + 1)
		html = html + "<link rel=\"next\" href=\"" + url1 + "\" />\n"
	}

	return template.HTML(html)
}

func (this *Pagination) Html() template.HTML {
	if this.total_items < this.per_page {
		return template.HTML("")
	}

	var html *bytes.Buffer = bytes.NewBufferString("<ul class=\"" + this.css + "\">")


	if this.First != "" {
		html.WriteString("<li class=\"first\"><a href=\"" + this.first_url + "\">" + this.First + "</a></li>")
	}

	page_range := this.Pages()

	if this.current > 1 {
		url1 := this.url + this.Prev2Current() //strconv.Itoa(this.current - 1)
		html.WriteString("<li class=\"prev\"><a rel=\"prev\" href=\"" + url1 + "\">" + this.Prev + "</a></li>")
	}

	//fmt.Println(page_range)
	for _, p := range page_range {
		n := strconv.Itoa(p)
		url1 := this.url + n
		class := ""
		if p == this.current {
			class = "active"
		}
		if (p == 1) {
			html.WriteString("<li class=\"" + class + "\"><a href=\"" + this.first_url + "\">" + n + "</a></li>")
		} else {
			html.WriteString("<li class=\"" + class + "\"><a href=\"" + url1 + "\">" + n + "</a></li>")
		}
	}

	if this.current < this.num_pages {
		url1 := this.url + this.Next2Current() //strconv.Itoa(this.current + 1)
		html.WriteString("<li class=\"next\"><a rel=\"next\" href=\"" + url1 + "\">" + this.Next + "</a></li>")
	}

	if this.Last != "" {
		url1 := this.url + strconv.Itoa(this.num_pages)
		html.WriteString("<li class=\"last\"><a href=\"" + url1 + "\">" + this.Last + "</a></li>")
	}

	html.WriteString("</ul>")
	return template.HTML(html.String())
}