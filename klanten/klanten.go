package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// Gegevens van de klanten die getoond, ingevoerd, gewijzigd of verwijderd kunnen worden
type Klant struct {
	Knr, Hnr, Ink                    int
	Nm, Vnm, Pc, Hnrt, Gsl, Blg, Rhf string
	Gbd, Krg, Opl, Opm               string
	Brf                              float32
}

type Bestelling struct {
	Bnr, Dlt, Knr, Vkp     int
	Sts, bsd, Klvnm, Klanm string
	Bdr                    float32
	afb					   decimal
}

// Verbinding maken met de database
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "vitaintellectdb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// De routing die afgeleid wordt naar de template
var tmpl = template.Must(template.ParseGlob("klanten/form/*"))

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

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
	err := db.QueryRow("SELECT voorletters, datum_in_dienst FROM medewerker WHERE voorletters=? AND datum_in_dienst=?", username, pass).Scan(&databaseUsername, &databasePassword)

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
<a href="/show">Mijn klanten</a> |
<a href="/index">Mijn bestellingen</a> |
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

/**
	KLANTEN
 */

//Index van alle klanten
func Index(w http.ResponseWriter, request *http.Request) {
	db := dbConn()
	userName := getUserName(request)
	selDB, err := db.Query("SELECT k.* FROM klant AS k INNER JOIN bestelling AS b ON b.klantnummer = k.klantnummer WHERE  b.verkoper = (SELECT m.medewerkernummer FROM medewerker AS m WHERE m.voorletters = ?)", userName) // Selecteren en ordenen van de gegevens van de klanten
	if err != nil {
		panic(err.Error())
	}
	bstl := Bestelling{}
	res := []Bestelling{}
	for selDB.Next() {
		var klantnummer, huisnummer, inkomen int
		var naam, voornaam, postcode, huisnummer_toevoeging, geslacht, bloedgroep, rhesusfactor string
		var geboortedatum, kredietregistratie, opleiding, opmerkingen string
		var beroepsrisicofactor float32
		err = selDB.Scan(&klantnummer, &voornaam, &naam, &postcode, &huisnummer, &huisnummer_toevoeging, &geboortedatum, &geslacht, &bloedgroep, &rhesusfactor, &beroepsrisicofactor, &inkomen, &kredietregistratie, &opleiding, &opmerkingen)
		if err != nil {
			panic(err.Error())
		}
		klnt.Knr = klantnummer
		klnt.Nm = naam
		klnt.Vnm = voornaam
		klnt.Pc = postcode
		klnt.Hnr = huisnummer
		klnt.Hnrt = huisnummer_toevoeging
		klnt.Gbd = geboortedatum
		klnt.Gsl = geslacht
		klnt.Blg = bloedgroep
		klnt.Rhf = rhesusfactor
		klnt.Brf = beroepsrisicofactor
		klnt.Ink = inkomen
		klnt.Krg = kredietregistratie
		klnt.Opl = opleiding
		klnt.Opm = opmerkingen
		res = append(res, klnt)
	}
	tmpl.ExecuteTemplate(w, "Index", nil)
	//tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

// De gegevens van één klant tonen
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nKnr := r.URL.Query().Get("klantnummer")
	selDB, err := db.Query("SELECT * FROM klant WHERE klantnummer=?", nKnr)
	if err != nil {
		panic(err.Error())
	}
	klnt := Klant{}
	for selDB.Next() {
		var klantnummer, huisnummer, inkomen int
		var naam, voornaam, postcode, huisnummer_toevoeging, geslacht, bloedgroep, rhesusfactor string
		var geboortedatum, kredietregistratie, opleiding, opmerkingen string
		var beroepsrisicofactor float32
		err = selDB.Scan(&klantnummer, &voornaam, &naam, &postcode, &huisnummer, &huisnummer_toevoeging, &geboortedatum, &geslacht, &bloedgroep, &rhesusfactor, &beroepsrisicofactor, &inkomen, &kredietregistratie, &opleiding, &opmerkingen)
		if err != nil {
			panic(err.Error())
		}
		klnt.Knr = klantnummer
		klnt.Nm = naam
		klnt.Vnm = voornaam
		klnt.Pc = postcode
		klnt.Hnr = huisnummer
		klnt.Hnrt = huisnummer_toevoeging
		klnt.Gbd = geboortedatum
		klnt.Gsl = geslacht
		klnt.Blg = bloedgroep
		klnt.Rhf = rhesusfactor
		klnt.Brf = beroepsrisicofactor
		klnt.Ink = inkomen
		klnt.Krg = kredietregistratie
		klnt.Opl = opleiding
		klnt.Opm = opmerkingen
	}
	tmpl.ExecuteTemplate(w, "Show", klnt)
	defer db.Close()
}

// Nieuwe klant toevoegen
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

//  Gegevens klant toevoegen aan de database
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		klantnummer := r.FormValue("klantnummer")
		voornaam := r.FormValue("voornaam")
		naam := r.FormValue("naam")
		postcode := r.FormValue("postcode")
		huisnummer := r.FormValue("huisnummer")
		huisnummer_toevoeging := r.FormValue("huisnummer_toevoeging")
		geboortedatum := r.FormValue("geboortedatum")
		geslacht := r.FormValue("geslacht")
		bloedgroep := r.FormValue("bloedgroep")
		rhesusfactor := r.FormValue("rhesusfactor")
		beroepsrisicofactor := r.FormValue("beroepsrisicofactor")
		inkomen := r.FormValue("inkomen")
		kredietregistratie := r.FormValue("kredietregistratie")
		opleiding := r.FormValue("opleiding")
		opmerkingen := r.FormValue("opmerkingen")
		insForm, err := db.Prepare("INSERT INTO klant(klantnummer, naam, voornaam, postcode, huisnummer, huisnummer_toevoeging, geboortedatum, geslacht, bloedgroep, rhesusfactor, beroepsrisicofactor, inkomen, kredietregistratie, opleiding, opmerkingen) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(klantnummer, naam, voornaam, postcode, huisnummer, huisnummer_toevoeging, geboortedatum, geslacht, bloedgroep, rhesusfactor, beroepsrisicofactor, inkomen, kredietregistratie, opleiding, opmerkingen)
		log.Println("Toegevoegd: Voornaam: " + voornaam + " | Achternaam: " + naam + " | Geboortedatum: " + geboortedatum)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Klant verwijderen uit de database
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	std := r.URL.Query().Get("klantnummer")
	delForm, err := db.Prepare("DELETE FROM klant WHERE klantnummer=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(std)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

/**
	BESTELLINGEN
 */

//Index van alle bestellingen
func bestellingIndex(w http.ResponseWriter, request *http.Request) {
	db := dbConn()
	userName := getUserName(request)
	selDB, err := db.Query("SELECT b.* FROM bestelling AS b WHERE b.verkoper =  (SELECT m.medewerkernummer FROM medewerker AS m WHERE m.voorletters = ? )", userName) // Selecteren en ordenen van de gegevens van de klanten
	if err != nil {
		panic(err.Error())
	}
	klnt := Klant{}
	res := []Klant{}
	for selDB.Next() {
		var bestelnummer, afbetaling_doorlooptijd, klantnummer, verkoper int
		var status, besteldatum, postcode, huisnummer_toevoeging, geslacht, bloedgroep, rhesusfactor string
		var geboortedatum, kredietregistratie, opleiding, opmerkingen string
		var beroepsrisicofactor float32
		err = selDB.Scan(&klantnummer, &voornaam, &naam, &postcode, &huisnummer, &huisnummer_toevoeging, &geboortedatum, &geslacht, &bloedgroep, &rhesusfactor, &beroepsrisicofactor, &inkomen, &kredietregistratie, &opleiding, &opmerkingen)
		if err != nil {
			panic(err.Error())
		}
		klnt.Knr = klantnummer
		klnt.Nm = naam
		klnt.Vnm = voornaam
		klnt.Pc = postcode
		klnt.Hnr = huisnummer
		klnt.Hnrt = huisnummer_toevoeging
		klnt.Gbd = geboortedatum
		klnt.Gsl = geslacht
		klnt.Blg = bloedgroep
		klnt.Rhf = rhesusfactor
		klnt.Brf = beroepsrisicofactor
		klnt.Ink = inkomen
		klnt.Krg = kredietregistratie
		klnt.Opl = opleiding
		klnt.Opm = opmerkingen
		res = append(res, klnt)
	}
	tmpl.ExecuteTemplate(w, "Index", nil)
	//tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

// De gegevens tonen
var router = mux.NewRouter()

func main() {
	log.Println("Server started on: http://localhost:8080")
	// login
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/internal", internalPageHandler)
	http.Handle("/", router)
	// klanten
	http.HandleFunc("/index", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	// bestellingen
	http.HandleFunc("/bestelling-index", bestellingIndex)
	http.HandleFunc("/bestelling-show", bestellingShow)
	http.HandleFunc("/bestelling-new", bestellingNew)
	http.HandleFunc("/bestelling-insert", bestellingInsert)
	http.HandleFunc("/bestelling-delete", bestellingDelete)
	// Start server
	http.ListenAndServe(":8080", nil)
}
