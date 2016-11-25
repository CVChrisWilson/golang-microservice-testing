package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
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
  id := insertToSql()
  selectFromSql(id)
}

func insertToSql() (int64) {
  db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/database_name?charset=utf8")
  checkErr(err)

  // insert
  stmt, err := db.Prepare("INSERT INTO portfolio SET created_by=?,updated_by=?,name=?,description=?,text_html=?")
  checkErr(err)

  res, err := stmt.Exec("Chris Wilson", "Chris Wilson", "API_Test", "Testing portfolio entry from API", "<h1>API Test INSERT</h1>")
  checkErr(err)

  id, err := res.LastInsertId()
  checkErr(err)

  return id
}

func selectFromSql(id int64) {
  db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/database_name?charset=utf8")
  checkErr(err)

  // select
  fmt.Println("SELECT * FROM portfolio WHERE id = '", id, "'")
  rows, err := db.Query("SELECT * FROM portfolio")
  checkErr(err)

  for rows.Next() {
    var id int
    var created_time string
    var updated_time string
    var created_by string
    var updated_by string
    var name string
    var description string
    var text_html string
    var demo_url sql.NullString
    var author sql.NullString
    err = rows.Scan(&id, &created_time, &updated_time, &created_by, &updated_by, &name, &description, &text_html, &demo_url, &author)
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
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
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
