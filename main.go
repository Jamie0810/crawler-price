package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anaskhan96/soup"
	"github.com/gorilla/mux"
)

type Product struct {
	Name  string
	Price string
}

type myJSON struct {
	Array []Product
}

var input []Product

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", getData).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getData(w http.ResponseWriter, r *http.Request) {
	ch := make(chan string)
	go parseUrls("https://tw.buy.yahoo.com/search/product?p=iphone", ch)
	data := <-ch
	fmt.Fprintf(w, data)
}

func parseUrls(url string, ch chan string) {

	doc := fetch(url)
	content := doc.Find("div", "class", "main").FindAll("span", "class", "BaseGridItem__itemInfo___3E5Bx")

	for i := 0; i < len(content); i += 1 {
		prodName := content[i].Find("span", "class", "BaseGridItem__title___2HWui").Text()
		prodPrice := content[i].Find("em", "class", "BaseGridItem__price___31jkj").Text()
		input = append(input, Product{prodName, prodPrice})
	}
	jsonData, _ := json.Marshal(&myJSON{input})

	ch <- string(jsonData)
}

func fetch(url string) soup.Root {
	soup.Headers = map[string]string{
		"User-Agent": "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	}

	resp, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	doc := soup.HTMLParse(resp)
	return doc
}
