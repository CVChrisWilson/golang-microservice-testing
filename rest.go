package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Payload struct {
    Portfolios Data
}

type Data struct {
    Portfolio Portfolios
}

type Portfolios struct {
    PortfolioItem PortfolioItems
    PortfolioTag PortfolioTags
}

type PortfolioItems map[string]int
type PortfolioTags map[string]int

func serveRest(w http.ResponseWriter, r *http.Request) {
	response, err := getJsonResponse()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(response))
}

func main() {
	http.HandleFunc("/", serveRest)
	http.ListenAndServe("localhost:1337", nil)
}

func getJsonResponse() ([]byte, error) {
	portfolioitems:= make(map[string]int)
	portfolioitems["SomeProject"] = 0
	portfolioitems["SomeOtherProject"] = 1
    portfolioitems["AnotherProject"] = 2

	portfoliotags := make(map[string]int)
	portfoliotags["C++"] = 21
	portfoliotags["Java"] = 0
    portfoliotags["Copyright Law"] = 62

	portfolios := Portfolios{portfolioitems, portfoliotags}
    d := Data{portfolios}
	p := Payload{d}

	return json.MarshalIndent(p, "", "  ")
}
