package main

import (
	"database/sql"
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

type Verkoper struct {
	Mnr  int
	Vanm string
	Vvnm string
}

// Verbinding maken met de database
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

// De routing die afgeled wordt naar de templat
var tmpl = template.Must(template.ParseGlob("form/*"))

// Index van alle klanten
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM klant ORDER BY klantnummer ASC") // Selecteren en ordenen van de gegevens van de klanten
	if err != nil {
		panic(err.Error())
	}
	klnt := Klant{}
	res := []Klant{}
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
	tmpl.ExecuteTemplate(w, "Index", res)
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

// De gegevens van een bestaande klant aanpassen
func Edit(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "Edit", klnt)
	defer db.Close()
}

// De gegevens toevoegen aan de database
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

// De gegevens veranderen in de database
func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
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
		klantnummer := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE klant SET voornaam=?, naam=?, postcode=?, huisnummer=?, huisnummer_toevoeging=?, geboortedatum=?, geslacht=?, bloedgroep=?, rhesusfactor=?, beroepsrisicofactor=?, inkomen=?, kredietregistratie=?, opleiding=?, opmerkingen=? WHERE klantnummer=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(naam, voornaam, postcode, huisnummer, huisnummer_toevoeging, geboortedatum, geslacht, bloedgroep, rhesusfactor, beroepsrisicofactor, inkomen, kredietregistratie, opleiding, opmerkingen, klantnummer)
		log.Println("Aangepast: Voornaam: " + voornaam + " | Achternaam: " + naam + " | Geboortedatum: " + geboortedatum)
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

// De gegevens tonen
func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
