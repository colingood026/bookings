package render

import (
	"bytes"
	"github.com/colingood026/bookings/internal/config"
	"github.com/colingood026/bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// RenderTemplate renders a template, 每次 request 過來都會從硬碟讀取一次 ，不是很有效率
//func RenderTemplate(w http.ResponseWriter, tmpl string) {
//	log.Println("execute RenderTemplate", tmpl)
//	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
//	err := parsedTemplate.Execute(w, nil)
//	if err != nil {
//		log.Println("error parsing template", err)
//		return
//	}
//
//}

//var templateCache = make(map[string]*template.Template)
//var lock sync.Mutex

// RenderTemplateFromCache renders a template, read from cache at first
//func RenderTemplateFromCache(w http.ResponseWriter, t string) {
//	log.Println("execute RenderTemplateFromCache", t)
//	var tmpl *template.Template
//	var err error
//	lock.Lock() // 遇到請求併發時的處理
//	_, inMap := templateCache[t]
//	// check if template already in templateCache
//	if !inMap {
//		log.Println("creating template and adding to cache", t)
//		err = createTemplateCache(t)
//		if err != nil {
//			log.Fatalln("error parse template:", t, err)
//		}
//	} else {
//		log.Println("load template from cache", t)
//	}
//	lock.Unlock()
//
//	tmpl = templateCache[t]
//
//	tmpl.Execute(w, nil)
//}

//func createTemplateCache(t string) error {
//	templates := []string{
//		fmt.Sprintf("./templates/%s", t),
//		"./templates/base.layout.tmpl",
//	}
//
//	parsedTemplate, err := template.ParseFiles(templates...)
//	if err != nil {
//		return err
//	}
//
//	templateCache[t] = parsedTemplate
//	return nil
//
//}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplateFromCacheV2(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	log.Println("execute RenderTemplateFromCacheV2", tmpl)
	var tc map[string]*template.Template
	if !app.UseCache {
		// only in dev mode !!!
		tc, _ = CreateTemplateCacheV2()
	} else {
		tc = app.TemplateCache
	}
	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("template not exist ", tmpl)
	}
	// 不一定要做，先做檢查
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}
	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
	}
}

func CreateTemplateCacheV2() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ./templates/
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	hasLayout, err := isLayoutExists()
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		if hasLayout {
			tmpl, err = tmpl.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[tmpl.Name()] = tmpl
	}
	return myCache, nil

}

func isLayoutExists() (bool, error) {
	matches, err := filepath.Glob("./templates/*.layout.tmpl")
	if err != nil {
		return false, err
	}
	return len(matches) > 0, nil
}
