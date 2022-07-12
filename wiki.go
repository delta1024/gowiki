package main

import (
	"os"
	"log"
	"regexp"
	"net/http"
	"errors"
	"html/template"
)
var tmpl_dir = "tmpl/"
var data_dir = "data/"
var templates = template.Must(template.ParseFiles(tmpl_dir+"edit.html", tmpl_dir+"view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var valid_link = regexp.MustCompile(`(\[)([a-zA-z0-9]+)(\])`)
var valid_new_line = regexp.MustCompile(`\n`)

type Page struct {
	Title string
	Body template.HTML
}

func (p *Page) save() error {
	filename := data_dir + p.Title + ".txt"
	return os.WriteFile(filename, []byte(p.Body), 0600)
}

func loadPage(title string) (*Page, error) {
	filename := data_dir +  title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	body_s := template.HTMLEscapeString(string(body))
	return &Page{Title: title, Body: template.HTML(body_s)}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {

	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	result := string(valid_link.ReplaceAllString(string(p.Body), "<a href=\"/view/$2\">$2</a>"))
	result = string(valid_new_line.ReplaceAllString(result, "<br/>"))
	p.Body = template.HTML(result)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)

}

func savedHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: template.HTML(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	title := "FrontPage"
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}



func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(savedHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
