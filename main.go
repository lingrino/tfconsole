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

type d struct {
	Data string
}

func (d *d) handler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("eval")
	d.Data = evalCombined(body)

	t, err := template.ParseFiles("templates/root.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, d)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	d := d{Data: ""}
	http.HandleFunc("/", d.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func evalCombined(line string) string {
	sOut, sErr, _ := eval(line)
	return fmt.Sprintf("%s%s", sOut, sErr)
}

// eval takes a line and runs it through terraform console.
// returns stdin, stderr, and the exit code.
func eval(line string) (string, string, int) {
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
