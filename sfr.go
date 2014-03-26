package main

import (
	"encoding/xml"
	"github.com/codegangsta/martini"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Item struct {
	Title       string        `xml:"title"`
	Link        template.HTML `xml:"link"`
	Description template.HTML `xml:"description"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

func main() {
	m := martini.Classic()
	m.Get("/", handler)
	m.Run()
}

func handler(w http.ResponseWriter, l *log.Logger) {
	feeds, err := getFeeds(l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t := template.Must(template.ParseFiles("templates/feeds.html"))
	err = t.Execute(w, feeds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getFeeds(l *log.Logger) (*Channel, error) {
	var rss RSS
	resp, err := http.Get("http://www.baixaki.com.br/rss/tecnologia.xml")
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	l.Printf("Content = %s", contents)
	err = xml.Unmarshal([]byte(contents), &rss)
	if err != nil {
		return nil, err
	}
	l.Printf("RSS = %#v", rss)
	return &rss.Channel, nil
}
