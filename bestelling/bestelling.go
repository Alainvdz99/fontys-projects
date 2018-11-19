package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// Gegevens van de bestelling die getoond, ingevoerd, gewijzigd of verwijderd kunnen worden
type Bestelling struct {
	Bnr, Dlt, Knr, Vkp     int
	Sts, bsd, Klvnm, Klanm string
	Bdr                    float32
}

// Gegevens van de klant, ter controle
type Klanten struct {
	Knr               int
	Klvnm, Klanm, Gbd string
}

// Gegevens van de verkoper ter koppeling tussen klant en bestelling
type Verkopers struct {
	Mnr          int
	Vkanm, Vkvnm string
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

// Het pad naar de templates
var tmpl = template.Must(template.ParseGlob("form/*"))

// Index van alle bestellingen
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	// Selecteren en ordenen van de gegevens van de bestellingen
	selDB, err := db.Query("SELECT b.bestelnummer, b.status, b.besteldatum, b.afbetaling_doorlooptijd, b.afbetaling_maandbedrag, k.klantnummer, k.voornaam, k.naam, m.medewerkernummer, m.voorletters, m.achternaam FROM bestelling b JOIN klant k ON b.klantnummer = k.klantnummer FROM bestelling b JION medewerker m ON b.verkoper = m.mederwerkernummer ORDER BY bestelnummer ASC")
	if err != nil {
		panic(err.Error())
	}
	bst := Bestelling{}
	res := []Bestelling{}
	for selDB.Next() {
		var bestelnummer, afbetaling_doorlooptijd int
		var status, voornaam, naam, voorletters, achternaam string
		var afbetaling_maandbedrag float32
		err = selDB.Scan(&bst.Bnr, &Bst.Sts, &Bst.Dlt, &Bst.Bdr, &Bst.Klvnm, &Bst.Klanm, &Bst.V)
		if err != nil {
			panic(err.Error())
		}
		B.Nr = nummer
		std.Vnm = voornaam
		std.Anm = achternaam
		std.Kls = klas
		res = append(res, std)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

// De gegevens van één student tonen
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nNr := r.URL.Query().Get("nummer")
	selDB, err := db.Query("SELECT * FROM student WHERE nummer=?", nNr)
	if err != nil {
		panic(err.Error())
	}
	std := Student{}
	for selDB.Next() {
		var nummer int
		var voornaam, achternaam, klas string
		err = selDB.Scan(&nummer, &voornaam, &achternaam, &klas)
		if err != nil {
			panic(err.Error())
		}
		std.Nr = nummer
		std.Vnm = voornaam
		std.Anm = achternaam
		std.Kls = klas
	}
	tmpl.ExecuteTemplate(w, "Show", std)
	defer db.Close()
}

// Nieuwe student toevoegen
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// De gegevens van een bestaande student aanpassen
func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nNr := r.URL.Query().Get("nummer")
	selDB, err := db.Query("SELECT * FROM student WHERE nummer=?", nNr)
	if err != nil {
		panic(err.Error())
	}
	std := Student{}
	for selDB.Next() {
		var nummer int
		var voornaam, achternaam, klas string
		err = selDB.Scan(&nummer, &voornaam, &achternaam, &klas)
		if err != nil {
			panic(err.Error())
		}
		std.Nr = nummer
		std.Vnm = voornaam
		std.Anm = achternaam
		std.Kls = klas
	}
	tmpl.ExecuteTemplate(w, "Edit", std)
	defer db.Close()
}

// De gegevens toevoegen aan de database
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		nummer := r.FormValue("nummer")
		voornaam := r.FormValue("voornaam")
		achternaam := r.FormValue("achternaam")
		klas := r.FormValue("klas")
		insForm, err := db.Prepare("INSERT INTO student(nummer, voornaam, achternaam, klas) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(nummer, voornaam, achternaam, klas)
		log.Println("Toegevoegd: Voornaam: " + voornaam + " | Achternaam: " + achternaam + " | Klas: " + klas)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// De gegevens veranderen in de database
func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		voornaam := r.FormValue("voornaam")
		achternaam := r.FormValue("achternaam")
		klas := r.FormValue("klas")
		nummer := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE student SET voornaam=?, achternaam=?, klas=? WHERE nummer=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(voornaam, achternaam, klas, nummer)
		log.Println("Aangepast: Voornaam: " + voornaam + " | Achternaam: " + achternaam + " | Klas: " + klas)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Student verwijderen in de database
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	std := r.URL.Query().Get("nummer")
	delForm, err := db.Prepare("DELETE FROM student WHERE nummer=?")
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
