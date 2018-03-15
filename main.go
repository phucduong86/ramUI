package main

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	//"github.com/gorilla/schema"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

////////////////////////////////////////
//logging middleware////////////////////
////////////////////////////////////////

func myLoggingHandler(h http.Handler) http.Handler {
	logFile, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return handlers.LoggingHandler(logFile, h)
}

var Store = sessions.NewCookieStore([]byte("suP3RaM"))

var templates map[string]*template.Template

func init() {
	loadTemplates()
}
func loadTemplates() {
	var baseTemplate = "templates/layout/_base.html"
	templates = make(map[string]*template.Template)

	templates["login"] = template.Must(template.ParseFiles(baseTemplate, "templates/login.html"))
	templates["userview"] = template.Must(template.ParseFiles(baseTemplate, "templates/userview.html"))
	templates["useradd"] = template.Must(template.ParseFiles(baseTemplate, "templates/useradd.html"))
	templates["useredit"] = template.Must(template.ParseFiles(baseTemplate, "templates/useredit.html"))

	templates["facilityview"] = template.Must(template.ParseFiles(baseTemplate, "templates/facilityview.html"))
	templates["facilityadd"] = template.Must(template.ParseFiles(baseTemplate, "templates/facilityadd.html"))
	templates["facilitydit"] = template.Must(template.ParseFiles(baseTemplate, "templates/facilityedit.html"))

	templates["formview"] = template.Must(template.ParseFiles(baseTemplate, "templates/formview.html"))
	templates["formadd"] = template.Must(template.ParseFiles(baseTemplate, "templates/formadd.html"))
	templates["formedit"] = template.Must(template.ParseFiles(baseTemplate, "templates/formedit.html"))
}
func main() {

	loginHandler := http.HandlerFunc(loginRoute)

	userAddHandler := http.HandlerFunc(userAddRoute)
	userViewHandler := http.HandlerFunc(userViewRoute)
	userEditHandler := http.HandlerFunc(userEditRoute)

	formViewHandler := http.HandlerFunc(indexRoute)
	formAddHandler := http.HandlerFunc(formAddRoute)
	formEditHandler := http.HandlerFunc(formEditRoute)

	facilityAddHandler := http.HandlerFunc(facilityAddRoute)
	facilityViewHandler := http.HandlerFunc(facilityViewRoute)
	facilityEditHandler := http.HandlerFunc(facilityEditRoute)

	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	router.PathPrefix("/user/static").Handler(http.StripPrefix("/user/static/", http.FileServer(http.Dir("static/"))))
	router.PathPrefix("/facility/static").Handler(http.StripPrefix("/facility/static/", http.FileServer(http.Dir("static/"))))
	router.PathPrefix("/form/static").Handler(http.StripPrefix("/form/static/", http.FileServer(http.Dir("static/"))))

	router.Handle("/login", myLoggingHandler(loginHandler))

	router.Handle("/", myLoggingHandler(formViewHandler)).Methods("GET")

	router.Handle("/user", myLoggingHandler(userViewHandler)).Methods("GET")
	router.Handle("/user/add", myLoggingHandler(userAddHandler)).Methods("GET")
	router.Handle("/user/edit", myLoggingHandler(userEditHandler)).Methods("GET")

	router.Handle("/form", myLoggingHandler(formViewHandler)).Methods("GET")
	router.Handle("/form/add", myLoggingHandler(formAddHandler)).Methods("GET")
	router.Handle("/form/edit", myLoggingHandler(formEditHandler)).Methods("GET")

	router.Handle("/facility", myLoggingHandler(facilityViewHandler)).Methods("GET")
	router.Handle("/facility/add", myLoggingHandler(facilityAddHandler)).Methods("GET")
	router.Handle("/facility/edit", myLoggingHandler(facilityEditHandler)).Methods("GET")

	if err := http.ListenAndServe(":4000", router); err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Println("Error fetching session cookie")
	}

	userRole := session.Values["userRole"]
	isLoggedIn := session.Values["loggedin"]

	m := map[string]interface{}{
		"role":     userRole,
		"loggedin": isLoggedIn,
	}

	if isLoggedIn == "true" {
		if err := templates["formview"].Execute(w, m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		log.Println("GET /: USER NOT LOGGED IN YET")
		http.Redirect(w, r, "/login", 302)
	}

}

func loginRoute(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Println("FAILED TO FETCH SESSION")
		if err := templates["login"].Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		isLoggedIn := session.Values["loggedin"] //after fetching sessions, find isLoggedIn value

		if isLoggedIn != "true" {
			if r.Method == "POST" {
				log.Println("POST /login then Pretent I successfully logged in")
				//validate user/password here then set the session cookie below
				session.Values["loggedin"] = "true"
				session.Values["userRole"] = "admin"
				session.Save(r, w)
				http.Redirect(w, r, "/", 302)
				log.Println("Saved a session cookie")
				return
			} else if r.Method == "GET" {
				log.Println("GET /login: Cookie indicate user not logged in when user get the page") //then render login page
				if err := templates["login"].Execute(w, nil); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

		} else {
			log.Println("Found a sweet cookie, redirect to index page")
			http.Redirect(w, r, "/", 302) //redirect to index page if cookie validated
		}
	}

}

func userAddRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["useradd"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func userViewRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["userview"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func userEditRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["useredit"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func facilityAddRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["facilityadd"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func facilityViewRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["facilityview"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func facilityEditRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["facilityedit"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func formAddRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["formadd"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func formEditRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["formedit"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
