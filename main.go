package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/DTherHtun/hola-unhealthy/statik"
	"github.com/rakyll/statik/fs"
)

//Info is info
type Info struct {
	Title string
	Done  string
}

//InfoPageData is data
type InfoPageData struct {
	PageTitle string
	Infos     []Info
}

var requests int64 = 0

func incRequests() int64 {
	return atomic.AddInt64(&requests, 1)
}

func getRequests() int64 {
	return atomic.LoadInt64(&requests)
}

func main() {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(statikFS)
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tmpl := template.Must(template.ParseFiles("/go/bin/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		incRequests()
		if getRequests() > 5 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			log.Printf("--> %s %s 500", r.Method, r.URL.Path)
		} else {
			data := InfoPageData{
				PageTitle: "Welcome! Hola",
				Infos: []Info{
					{Title: "Pod Name - ", Done: hostName},
					{Title: "Host - ", Done: r.Host},
					{Title: "Remote Address - ", Done: r.RemoteAddr},
					{Title: "Request Method - ", Done: r.Method},
					{Title: "Request URL - ", Done: r.URL.String()},
					{Title: "Protocol - ", Done: r.Proto},
					{Title: "UserAgent - ", Done: r.UserAgent()},
				},
			}
			tmpl.Execute(w, data)
			log.Printf("--> %s %s 200", r.Method, r.URL.Path)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
