package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/models"
	"github.com/justinas/nosurf"
)

//allows us specify certain functions available to golang template
var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

//returns time in YYYY-MM-DD format in the templates
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

//adding 2 int tohether
func Add(a, b int) int {
	return a + b
}

//this func allows us to iterate within 2 days
//iterate returns a slice of int , starting at 1,going to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

var app *config.AppConfig

//this is important to get the templates when runnng tests
var pathToTemplates = "./templates"

//sets the configuration for a new template ,this helps to optimize our template cache which is stored in a map(TemplateCache)
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// these items will be automatically populated when a page is rendered
	td.Flash = app.Session.PopString(r.Context(), "flash") // PopString puts something in the session until the next time a page is displayed
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticted = 1 //user is loged in
	}
	return td
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	//using map to hold a particular data structure,creates 2 entries which are in about.page.html & home.page.html
	myCache := map[string]*template.Template{}
	//getting all the pages in the template folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	//blank identifiers are used for identfication for a particular purpose without having to refer or return it golang has a feature to refer and use it
	for _, page := range pages { //blank identifiers  to identify the pages
		name := filepath.Base(page) //extracts the name of the page about.page.html & home.page.html using filepath
		//creating a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		//getting layouts formats from the templates
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err //returns eror
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts //taking the template set and adding it to the myCache map,adding template to template cache
	}

	return myCache, nil
}

// Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, html string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	//tc which stores template cache
	//template cache stores new templates for retrival
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache //gets template cache from the appconfig(config dir)
	} else {
		// this is just used for testing, so that we rebuild
		// the cache on every request
		tc, _ = CreateTemplateCache()
	}
	//gets template
	t, ok := tc[html] //if the template exist it will have a value and ok =true
	//if template doesnt exist it will have no value ok =false
	if !ok {
		return errors.New("can't get template from cache")
	}
	//creating a buffer for a template that is not in the template dir or disk
	buf := new(bytes.Buffer) //puts the parsed template that is currently in memory into some bytes

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td) //takes the template  executes dont parse any data and store in the buf variable
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.WriteTo(w) //writing the response to the response writer
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}
