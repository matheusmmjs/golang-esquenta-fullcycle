package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

var courses []Course

type Course struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func generateCourses() {
	course1 := Course{
		ID:   "1",
		Name: "Full Cycle",
	}

	course2 := Course{
		ID:   "2",
		Name: "Bonus Full Cycle",
	}

	courses = append(courses, course1, course2)
}

func main() {
	generateCourses()

	e := echo.New()

	e.GET("/courses", listCourses)
	e.POST("/courses", createCourse)

	e.Logger.Fatal(e.Start(":8081"))
}

func listCourses(c echo.Context) error {
	return c.JSON(http.StatusOK, courses)
}

func createCourse(c echo.Context) error {
	course := Course{}
	c.Bind(&course)

	err := persistCourse(course)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, course)
}

func persistCourse(course Course) error {
	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		return err
	}

	stmt, err := db.Prepare("insert into courses values ($1, $2)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(course.ID, course.Name)

	if err != nil {
		return err
	}

	return nil

}
