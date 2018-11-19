package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

// Verbinding maken met de database Vita Intellectus
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Kikker12"
	dbName := "vitaintellectdb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// Instellen van de cookie
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Oproepen van de gebruikersnaam voor de sessie
func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

//Sessie instellen
func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// Inlogpagina instellen
func loginHandler(response http.ResponseWriter, request *http.Request) {
	db := dbConn()

	//Gebruikersnaam en wachtwoord van de webpagina lezen
	username := request.FormValue("username")
	pass := request.FormValue("password")

	var databaseUsername string
	var databasePassword string

	//Gebruiker in de database selecteren
	err := db.QueryRow("SELECT voorletters, datum_in_dienst FROM medewerker WHERE voorletters=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(response, request, "/login", 301)
		return
	}

	redirectTarget := "/"
	if username != "" && pass != "" {
		// .. check credentials ..
		setSession(username, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// index page
const indexPage = `
<h1>Inloggen</h1>
<form method="post" action="/login">
    <label for="username">Gebruikersnaam</label>
    <input type="text" id="username" name="username" placeholder="gebruikersnaam">
    <label for="password">Wachtwoord</label>
    <input type="password" id="password" name="password" placeholder="wachtwoord">
    <button type="submit">Inloggen</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// internal page
const internalPage = `<html><body>
<h1>Venditor: Vita Intellectus</h1>

<hr>
<p>Hallo: %s </p>
<p>
<a href="/mijnklant">Mijn klanten</a> |
<a href="/">Mijn bestellingen</a> |
<a href="/">Klanten</a> |
<a href="/">Bestellingen</a> |
<a href="/">Modules</a>
</p>
<form method="post" action="/logout">
    <button type="submit">Uitloggen</button>
</form>
</body></html>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)

	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func mijnklantPageHandler(response http.ResponseWriter, request *http.Request) {
	db := dbConn()
	userName := getUserName(request)
	db.Query("SELECT k.* FROM klant AS k INNER JOIN bestelling AS b ON b.klantnummer = k.klantnummer WHERE  b.verkoper = (SELECT m.medewerkernummer FROM medewerker AS m WHERE m.naam = ?)", userName)

	if userName != "" {
		fmt.Fprintf(response, mijnklantPage, userName)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// Pagina Mijn klanten
const mijnklantPage = `<html><body>
{{if .userName }}
<h1>Venditor: Vita Intellectus</h1>

<hr>
<p>Mijn klanten van <b>%s </b></p>
<p>
<body>
<h1>Mijn klanten</h1>
    <table border="1">
      <thead>
      <tr>
        <td>Nr.</td>
        <td>Naam</td>
        <td>Postcode/huisnummer</td>
        <td>Geboortedatum</td>
        <td>Bekijk</td>
        <td>Wijzig</td>
      </tr>
       </thead>
       <tbody>
      <tr>
        <td> %s </td>
        <td> {{ .voornaam }} {{ .naam}} </td>
        <td> {{ .postcode }} {{ .huisnummer }} {{ .huisnummertoevoeging }} </td>
        <td> {{ .geboortedatum }} </td> 
        <td><a href="/show?klantnummer={{ .klantnummer }}">Bekijk</a></td>
        <td><a href="/edit?klantnummer={{ .klantnummer }}">Wijzig</a></td>
      </tr>
</p>
{{end}}
</body></html>

`

// server main method
var router = mux.NewRouter()

func main() {

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)
	router.HandleFunc("/mijnklant", mijnklantPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
