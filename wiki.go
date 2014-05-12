package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

//New Page to implement HTML Views with escaped HTML
//without changing bytes stored in TXT Files in Page Type
type HTMLPage struct {
	Title string
	Body  template.HTML
}

var template_dir = "templates"
var templates = template.Must(template.ParseFiles(template_dir+"/edit.html", template_dir+"/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-0]+)$")
var reg, _ = regexp.Compile(`\[([a-zA-Z0-0]+)\]`) // Throws unknown escape sequence when using double quotes

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile("wikis/"+filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile("wikis/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // title second subexpression
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *HTMLPage) {
	//When preloading all templates, dont prefix the folder name
	// Just use filename instead .
	// FIXME: What to do when template names are same accross folders ?
	//eg: /templates/users/edit.html & /templates/wikis/edit.html
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s\n", r.URL.Path)
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2]) // title second subexpression
	}
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

func htmlize(p *Page) *HTMLPage {
	safe := template.HTMLEscapeString(string(p.Body))
	HPage := &HTMLPage{Title: p.Title, Body: template.HTML(safe)}
	return HPage
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	//So GOLang give you nil errors if there is something wrong with templates
	// http: panic serving 127.0.0.1:58332: runtime error: invalid memory address or nil pointer dereference
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	//Internal linking using [page]
	safe := template.HTMLEscapeString(string(p.Body))
	unescaped := reg.ReplaceAllString(string(safe), "<a href=\"/view/$1\">$1</a>")
	HPage := &HTMLPage{Title: p.Title, Body: template.HTML(unescaped)}
	renderTemplate(w, "view", HPage)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	ht_p := htmlize(p)
	renderTemplate(w, "edit", ht_p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
