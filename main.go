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
	Id    int
	Name  string
	Price string
}

type myJSON struct {
	Array []Product
}

var input []Product

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", getPrice).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getPrice(w http.ResponseWriter, r *http.Request) {
	ch := make(chan string)
	go parseUrls("https://www.best-price.com/search/result/query/bed/", ch)
	data := <-ch
	fmt.Fprintf(w, data)
}

func parseUrls(url string, ch chan string) {
	doc := fetch(url)
	content := doc.Find("ul", "class", "sc-eHgmQL").FindAll("li", "class", "sc-iAyFgw")

	for i := 0; i < len(content); i += 1 {
		id := i
		prodName := content[i].Find("span", "class", "byHTMx").Text()
		prodPrice := content[i].Find("span", "class", "joktCJ").Text()
		input = append(input, Product{id, prodName, prodPrice})
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
