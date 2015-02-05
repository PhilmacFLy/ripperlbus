package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type mitfahrer struct {
	Passagiere [6]string
	Disabled   [6]bool
}

func statichandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func (m *mitfahrer) Save() error {
	b, err := json.MarshalIndent(&m, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile("mitfahrer.json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (m *mitfahrer) Load() error {
	body, err := ioutil.ReadFile("mitfahrer.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	return nil
}

func savehandler(w http.ResponseWriter, r *http.Request) {
	p := mitfahrer{}

	for i := range p.Passagiere {
		val := r.FormValue(strconv.Itoa(i))
		if p.Passagiere[i] == "" {
			p.Passagiere[i] = val
		}
		p.Disabled[i] = (val != "")
	}
	p.Save()
	http.Redirect(w, r, "/", http.StatusFound)
}

func mainhandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template.html")

	if err != nil {
		log.Fatal(err)
	}

	m := mitfahrer{}

	m.Load()

	err = t.Execute(w, &m)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	if _, err := os.Stat("mitfahrer.json"); os.IsNotExist(err) {
		fmt.Printf("no such file or directory")
		m := mitfahrer{}
		m.Save()
	}

	http.HandleFunc("/static/", statichandler)
	http.HandleFunc("/", mainhandler)
	http.HandleFunc("/submit", savehandler)
	http.ListenAndServe(":1313", nil)
}
