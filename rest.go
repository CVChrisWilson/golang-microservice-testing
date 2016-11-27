package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	_ "github.com/go-sql-driver/mysql"
)

type Payload struct {
	Portfolios Data
}

type Data struct {
	Portfolio Portfolios
}

type Portfolios struct {
	PortfolioItem PortfolioItems
	PortfolioTag  PortfolioTags
}

type PortfolioItems map[string]int
type PortfolioTags map[string]int

func serveRest(w http.ResponseWriter, r *http.Request) {
	var sqlResponse string
	switch r.Method {
	case "GET":
		// serve the resource
		sqlResponse = selectFromSql(path.Base(r.URL.Path))
	case "POST":
		// insert to db
		htmlData, err := ioutil.ReadAll(r.Body)
		if err == nil {
			insertToSql(string(htmlData))
		}
	}
	//response, err := getJsonResponse()

	fmt.Fprintf(w, string(sqlResponse))
}

func main() {
	http.HandleFunc("/portfolio/", serveRest)
	http.ListenAndServe("localhost:1401", nil)
	//id := insertToSql()
	//selectFromSql(id)
}

func insertToSql(json string) int64 {
	db, err := sql.Open("mysql", "root:password!@tcp(localhost:3306)/db?charset=utf8")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT INTO portfolio SET created_by=?,updated_by=?,name=?,description=?,text_html=?")
	checkErr(err)
	fmt.Println("Chris Wilson", "Chris Wilson", "API_Test", "Testing portfolio entry from API", json)
	res, err := stmt.Exec("Chris Wilson", "Chris Wilson", "API_Test", "Testing portfolio entry from API", json)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

func selectFromSql(id string) string {
	db, err := sql.Open("mysql", "root:password!@tcp(localhost:3306)/db?charset=utf8")
	checkErr(err)

	// select
	qs := "SELECT * FROM portfolio WHERE id = '" + id + "' LIMIT 1"
	fmt.Println(qs)
	//rows, err := db.Query("SELECT * FROM portfolio")
	row, err := db.Query(qs)
	checkErr(err)

	var created_time string
	var updated_time string
	var created_by string
	var updated_by string
	var name string
	var description string
	var text_html string
	var demo_url sql.NullString
	var author sql.NullString
	for row.Next() {
		err = row.Scan(&id, &created_time, &updated_time, &created_by, &updated_by, &name, &description, &text_html, &demo_url, &author)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(created_time)
		fmt.Println(updated_time)
		fmt.Println(created_by)
		fmt.Println(updated_by)
		fmt.Println(name)
		fmt.Println(description)
		fmt.Println(text_html)
		fmt.Println(demo_url)
		fmt.Println(author)

	}
	return created_time + updated_time + created_by + updated_by + name + description + text_html
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getJsonResponse() ([]byte, error) {
	portfolioitems := make(map[string]int)
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
