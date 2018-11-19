package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

var (
	dbGebruikersnaam = os.Getenv("S3_BUCKET")
	dbWachtwoord     = os.Getenv("SECRET_KEY")
)

func main() {
	student := GeefStudent(3)
	fmt.Println(student.Voornaam) // Dit is Marieke.

	klas := GeefKlas("C34")
	for _, studentInKlas := range klas.Studenten {
		fmt.Println(studentInKlas.Voornaam) // Karel, Stephan en Linda
	}
}

func GeefStudent(nummer int) Student {
	login := dbGebruikersnaam + ":" + dbWachtwoord
	db, err := sql.Open("mysql", login+"@/studentenadministratie")
	if err != nil {
		panic(err)
	}
	query := "SELECT nummer, voornaam, achternaam FROM student WHERE nummer = ?"
	resultaat, err := db.Query(query, nummer)
	if err != nil {
		panic(err)
	}
	var student Student
	if resultaat.Next() {
		err := resultaat.Scan(&student.Nummer, &student.Voornaam, &student.Achternaam)
		if err != nil {
			panic(err)
		}
	}
	resultaat.Close()
	return student
}

func GeefKlas(code string) Klas {
	login := dbGebruikersnaam + ":" + dbWachtwoord
	db, err := sql.Open("mysql", login+"@/studentenadministratie")
	if err != nil {
		panic(err)
	}
	query := "SELECT s.nummer, s.voornaam, s.achternaam, k.code, k.locatie FROM student s JOIN klas k ON s.klas = k.code WHERE k.code = ?"
	resultaat, err := db.Query(query, code)
	if err != nil {
		panic(err)
	}
	var klas Klas
	klas.Studenten = make([]Student, 0)
	for resultaat.Next() {
		var student Student
		err := resultaat.Scan(&student.Nummer, &student.Voornaam, &student.Achternaam, &klas.Code, &klas.Locatie)
		if err != nil {
			panic(err)
		}
		klas.Studenten = append(klas.Studenten, student)
	}
	resultaat.Close()
	return klas
}
