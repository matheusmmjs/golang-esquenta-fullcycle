package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

var courses []Course

type Course struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	e := echo.New()

	e.GET("/courses", listCourses)
	e.POST("/courses", createCourse)

	e.Logger.Fatal(e.Start(":8081"))
}

func listCourses(c echo.Context) error {
	db, err := sql.Open("sqlite3", "test.db")

	stmt, err := db.Prepare("SELECT * FROM courses")

	if err != nil {
		return err
	}

	rows, err := stmt.Query()

	if err == rows.Err() {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err != nil {
		return err
	}

	defer stmt.Close()

	course := Course{}

	for rows.Next() {
		err = rows.Scan(&course.ID, &course.Name)

		if err == rows.Err() {
			return c.JSON(http.StatusUnprocessableEntity, nil)
		}

		if err != nil {
			return err
		}

		courses = append(courses, course)
	}

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

	defer stmt.Close()

	return nil
}
