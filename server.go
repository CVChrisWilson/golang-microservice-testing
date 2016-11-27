package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

type Portfolio struct {
	id           int64  `db:"id" json:"id"`
	created_time string `db:"created_time" json:"created_time"`
	updated_time string `db:"updated_time" json:"updated_time"`
	created_by   string `db:"created_by" json:"created_by"`
	updated_by   string `db:"updated_by" json:"updated_by"`
	name         string `db:"name" json:"name"`
	description  string `db:"description" json:"description"`
	text_html    string `db:"text_html" json:"text_html"`
	demo_url     string `db:"demo_url" json:"demo_url"`
	author       string `db:"author" json:"author"`
}

var dbmap = initDb()

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:password!@/db")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Portfolio{}, "portfolio").SetKeys(true, "id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		v1.GET("/portfolios", GetPortfolios)
		v1.GET("/portfolios/:id", GetPortfolio)
		v1.POST("/portfolios", PostPortfolio)
		v1.PUT("/portfolios/:id", UpdatePortfolio)
		v1.DELETE("/portfolios/:id", DeletePortfolio)
		v1.OPTIONS("/portfolios", OptionsPortfolio)     // POST
		v1.OPTIONS("/portfolios/:id", OptionsPortfolio) // PUT, DELETE
	}

	r.Run(":8888")
}

func GetPortfolios(c *gin.Context) {
	var portfolios []Portfolio
	_, err := dbmap.Select(&portfolios, "SELECT * FROM portfolio")

	if err == nil {
		c.JSON(200, portfolios)
	} else {
		c.JSON(404, gin.H{"error": "no portfolio(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/portfolios
}

func GetPortfolio(c *gin.Context) {
	id := c.Params.ByName("id")
	var portfolio Portfolio
	err := dbmap.SelectOne(&portfolio, "SELECT * FROM portfolio WHERE id=? LIMIT 1", id)

	if err == nil {
		portfolio_id, _ := strconv.ParseInt(id, 0, 64)

		content := &Portfolio{
			id:           portfolio_id,
			created_time: portfolio.created_time,
			updated_time: portfolio.updated_time,
			created_by:   portfolio.created_by,
			updated_by:   portfolio.updated_by,
			name:         portfolio.name,
			description:  portfolio.description,
			text_html:    portfolio.text_html,
			demo_url:     portfolio.demo_url,
			author:       portfolio.author,
		}
		c.JSON(200, content)
	} else {
		c.JSON(404, gin.H{"error": "portfolio not found"})
	}

	// curl -i http://localhost:8080/api/v1/portfolios/1
}

func PostPortfolio(c *gin.Context) {
	var portfolio Portfolio
	x, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("%s", string(x))

	c.Bind(&portfolio)

	log.Println(portfolio)

	if portfolio.created_by == "" {
		portfolio.created_by = "Anonymous"
	}

	if portfolio.author == "" {
		portfolio.author = portfolio.created_by
	}
	log.Println(portfolio.name)
	if portfolio.created_by != "" && portfolio.name != "" && portfolio.text_html != "" {

		if insert, _ := dbmap.Exec(`INSERT INTO portfolio (created_by, name, description, text_html, demo_url, author) VALUES (?, ?, ?, ?, ?, ?)`, portfolio.created_by, portfolio.name, portfolio.description, portfolio.text_html, portfolio.demo_url, portfolio.author); insert != nil {
			portfolio_id, err := insert.LastInsertId()
			if err == nil {
				content := &Portfolio{
					id:           portfolio_id,
					created_time: portfolio.created_time,
					updated_time: portfolio.updated_time,
					created_by:   portfolio.created_by,
					updated_by:   portfolio.updated_by,
					name:         portfolio.name,
					description:  portfolio.description,
					text_html:    portfolio.text_html,
					demo_url:     portfolio.demo_url,
					author:       portfolio.author,
				}
				c.JSON(201, content)
			} else {
				checkErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(400, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/portfolios
}

func UpdatePortfolio(c *gin.Context) {
	id := c.Params.ByName("id")
	var portfolio Portfolio
	err := dbmap.SelectOne(&portfolio, "SELECT * FROM portfolio WHERE id=?", id)

	if err == nil {
		var json Portfolio
		c.Bind(&json)

		portfolio_id, _ := strconv.ParseInt(id, 0, 64)

		portfolio := Portfolio{
			id:           portfolio_id,
			created_time: portfolio.created_time,
			updated_time: portfolio.updated_time,
			created_by:   portfolio.created_by,
			updated_by:   portfolio.updated_by,
			name:         portfolio.name,
			description:  portfolio.description,
			text_html:    portfolio.text_html,
			demo_url:     portfolio.demo_url,
			author:       portfolio.author,
		}

		if portfolio.created_by == "" {
			portfolio.created_by = "Anonymous"
		}

		if portfolio.updated_by == "" {
			portfolio.updated_by = portfolio.created_by
		}

		if portfolio.author == "" {
			portfolio.author = portfolio.created_by
		}

		if portfolio.created_by != "" && portfolio.updated_by != "" && portfolio.name != "" && portfolio.text_html != "" {
			_, err = dbmap.Update(&portfolio)

			if err == nil {
				c.JSON(200, portfolio)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "portfolio not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}

func DeletePortfolio(c *gin.Context) {
	id := c.Params.ByName("id")

	var portfolio Portfolio
	err := dbmap.SelectOne(&portfolio, "SELECT * FROM portfolio WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&portfolio)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "portfolio not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/users/1
}

func OptionsPortfolio(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}
