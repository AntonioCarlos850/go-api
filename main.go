package main

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Route struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

var people []Person
var routes []Route

func define() error {
	routes = append(routes, Route{"/people", "GET"})
	routes = append(routes, Route{"/people", "POST"})
	return nil
}

func defaultRoutes(c echo.Context) error {
	return c.JSON(200, routes)
}

func main() {
	define()
	e := echo.New()
	e.GET("/", defaultRoutes)
	e.POST("/people", createPerson)
	e.GET("/people", getPeople)
	e.Logger.Fatal(e.Start(":8080"))
}

func getPeople(c echo.Context) error {
	return c.JSON(200, people)
}

func createPerson(c echo.Context) error {
	person := new(Person)

	if err := c.Bind(person); err != nil {
		return err
	}

	people = append(people, *person)
	savePerson(*person)
	return c.JSON(200, people)
}

func savePerson(person Person) error {
	db, err := sql.Open("sqlite3", "people.db")
	if err != nil {
		return err
	}
	defer db.Close()

	smtp, err := db.Prepare("INSERT INTO people (name, age) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	_, err = smtp.Exec(person.Name, person.Age)
	if err != nil {
		return err
	}

	return nil
}
