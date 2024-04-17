package main

import (
	"html/template"
	"io"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var levels = []string{"error", "warn", "info", "debug"}

type Log struct {
	Level string
	Time  string
	Body  string
}

// getRandomLog generates a random log entry
func NewLog() Log {
	// Generate a random log body
	body := generateRandomString(30) // Adjust length as needed

	// Generate a random log level
	level := levels[rand.Intn(len(levels))]

	// Generate a random time within the last 24 hours
	randomTime := randomTimeWithinLastDay()
	timeStr := randomTime.Format("2006-01-02T15:04:05Z07:00") // RFC3339 format

	return Log{
		Time:  timeStr,
		Body:  body,
		Level: level,
	}
}

// generateRandomString generates a random string of given length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// randomTimeWithinLastDay generates a random time within the last 24 hours
func randomTimeWithinLastDay() time.Time {
	now := time.Now()
	max := now.Unix()
	min := now.Add(-24 * time.Hour).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	logs := []Log{NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog()}
	// for i, log := range logs {
	// 	e.Logger.Info(fmt.Sprintf("Log %d: %s - %s %s\n", i, log.Time, log.Level, log.Body))
	// }

	e.Static("/static", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", map[string]interface{}{
			"Logs": logs,
		})
	})

	e.Logger.Fatal(e.Start(":42069"))
}
