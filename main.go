package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// app stores shared information for the service like cached templates.
type app struct {
	templates *template.Template
}

// main starts the application.
func main() {
	app := &app{
		templates: template.Must(template.ParseGlob("./templates/*")),
	}

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.handlerConsole)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handlerConsole renders the console and handles new requests.
func (a *app) handlerConsole(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("input")

	err := a.templates.ExecuteTemplate(w, "console.html.tmpl", struct{ Line string }{Line: consoleCombined(input)})
	if err != nil {
		http.Error(w, "ERROR: failed to render template", http.StatusInternalServerError)
		log.Println("ERROR: failed to render console template:", err)
	}
}

// consoleCombined runs console() and combines stdout and stderr.
func consoleCombined(line string) string {
	sOut, sErr, _ := console(line)
	return fmt.Sprintf("%s%s", sOut, sErr)
}

// console takes a line and runs it through terraform console.
// returns stdin, stderr, and the exit code.
func console(line string) (string, string, int) {
	var code int
	var sOut, sErr bytes.Buffer

	cmd := exec.Command("terraform", "console")
	cmd.Stdin = strings.NewReader(line)
	cmd.Stdout = &sOut
	cmd.Stderr = &sErr

	err := cmd.Run()
	if err != nil {
		code = 999 // signals a rare error that is not an ExitError
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			code = ee.ExitCode()
		}
	}

	return sOut.String(), sErr.String(), code
}
