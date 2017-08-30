package main

import (
	"regexp"
	"io/ioutil"
	"net/http"
	"html/template"
	"encoding/json"
)

var tmplPath = "go/main/tmpl/"
var dataPath = "go/main/data/"
var templates = template.Must(template.ParseFiles(tmplPath + "edit.html", tmplPath + "view.html", tmplPath + "build/index.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var replaceAll = regexp.MustCompile(".txt")

type Page struct {
	Title string
	Body []byte
	PageNames []string
}

func main() {
	//http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/api/json", foo)
	http.HandleFunc("/lol/", lol)

	fs := http.FileServer(http.Dir("go/main/tmpl/build/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func lol(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

type Profile struct {
	Name    string
	Hobbies []string
}

func foo(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Page) save() error {
	filename := dataPath + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := dataPath + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	lol, _ := ioutil.ReadDir(dataPath)
	pageNames := []string{}
	for _, f := range lol {
		pageNames = append(pageNames, replaceAll.ReplaceAllString(f.Name(), ""))
	}
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body, PageNames: pageNames}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
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