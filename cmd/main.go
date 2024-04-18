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
	t := &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	return t
}

var levels = []string{"error", "warn", "info", "debug"}

type Log struct {
	Level string
	Time  string
	Body  string
}

type FilterButtons struct {
	Level string
	Label string
}

func NewLog() Log {
	level := levels[rand.Intn(len(levels))]

	return NewLogWithLevel(level)
}

func NewLogWithLevel(level string) Log {
	body := generateRandomString(30)

	randomTime := randomTimeWithinLastDay()
	timeStr := randomTime.Format("2006-01-02T15:04:05Z07:00")

	return Log{
		Time:  timeStr,
		Body:  body,
		Level: level,
	}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

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

	// serve .css files from /css
	e.Static("/static", "css")

	logs := []Log{NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog(), NewLog()}

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", map[string]interface{}{
			"Buttons": []FilterButtons{
				{"error", "Error"},
				{"warn", "Warn"},
				{"info", "Info"},
				{"debug", "Debug"},
				{"", "All"},
			},
		})
	})

	e.GET("/logs", func(c echo.Context) error {
		return c.Render(200, "logs", map[string]interface{}{
			"Logs": logs,
		})
	})

	e.POST("/level", func(c echo.Context) error {
		level := c.FormValue("level")
		filteredLogs := filterLogsByLevel(logs, level)
		return c.Render(200, "logs", map[string]interface{}{
			"Logs": filteredLogs,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func filterLogsByLevel(logs []Log, levelFilter string) []Log {
	if levelFilter == "" {
		return logs
	}

	var filteredLogs []Log
	for _, log := range logs {
		if log.Level == levelFilter {
			filteredLogs = append(filteredLogs, log)
		}
	}
	return filteredLogs
}
