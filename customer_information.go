package main

import "fmt"

type Work struct {

	title, company string
	workHoursWeekly, yearsExperience int
}

type Address struct {

	street, city, country, postalCodeLetters string
	streetNumber, postalCodeNumber int
}

type Birthday struct {
	dayOfBirth, monthOfBirth, yearOfBirth int
}

type Person struct {

	firstName, lastName string
	age, length, weight int
	address Address
	work Work
	birthday Birthday
}

func main (){
	var p1 Person
	p1.firstName = "Alain"
	p1.lastName  = "van der Zanden"
	p1.age		 = 19
	p1.length 	 = 182
	p1.weight	 = 69
	p1.address 	 = Address{
		street:	"Iepebeek",
		city:	"Veldhoven",
		country:	"Netherlands",
		streetNumber:	10,
		postalCodeNumber: 5501,
		postalCodeLetters: "CX",
	}
	p1.work		 = Work{
		title:	"Back-end developer",
		company:	"Toppershops",
		workHoursWeekly:	32,
		yearsExperience:	1,
	}
	p1.birthday	 = Birthday{
		dayOfBirth:	    21,
		monthOfBirth: 	5,
		yearOfBirth:	1999,
	}
	fmt.Println("Name:", p1.firstName, p1.lastName)
	fmt.Println("Date of birth:", p1.birthday.dayOfBirth,"-", p1.birthday.monthOfBirth,"-", p1.birthday.yearOfBirth)
	fmt.Println("Length:", p1.length, "cm")
	fmt.Println("Weight:", p1.weight, "kg")
	fmt.Println("Address:", p1.address.street, p1.address.streetNumber)
	fmt.Println("Country:", p1.address.country)
	fmt.Println("Postal Code:", p1.address.postalCodeNumber, p1.address.postalCodeLetters)
	fmt.Println("Function:", p1.work.title, "at", p1.work.company)
	fmt.Println("Work time weekly:", p1.work.workHoursWeekly, "hours")
	fmt.Println("Experience:", p1.work.yearsExperience, "year")

}
















