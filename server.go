package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	gorp "gopkg.in/gorp.v2"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Portfolio struct {
	Id           int64  `db:"id" json:"id"`
	Created_time string `db:"createdtime" json:"createdtime"`
	Updated_time string `db:"updatedtime" json:"updatedtime"`
	Created_by   string `db:"createdby" json:"createdby"`
	Updated_by   string `db:"updatedby" json:"updatedby"`
	Name         string `db:"name" json:"name"`
	Description  string `db:"description" json:"description"`
	Text_html    string `db:"texthtml" json:"texthtml"`
	Demo_url     string `db:"demourl" json:"demourl"`
	Author       string `db:"author" json:"author"`
}

var dbmap = initDb()

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db?charset=utf8")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Portfolio{}, "Portfolio").SetKeys(true, "Id")
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
	_, err := dbmap.Select(&portfolios, "SELECT * FROM Portfolio")

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
	err := dbmap.SelectOne(&portfolio, "SELECT * FROM Portfolio WHERE id=? LIMIT 1", id)

	if err == nil {
		portfolio_id, _ := strconv.ParseInt(id, 0, 64)

		content := &Portfolio{
			Id:           portfolio_id,
			Created_time: portfolio.Created_time,
			Updated_time: portfolio.Updated_time,
			Created_by:   portfolio.Created_by,
			Updated_by:   portfolio.Updated_by,
			Name:         portfolio.Name,
			Description:  portfolio.Description,
			Text_html:    portfolio.Text_html,
			Demo_url:     portfolio.Demo_url,
			Author:       portfolio.Author,
		}
		c.JSON(200, content)
	} else {
		fmt.Println(err)
		c.JSON(404, gin.H{"error": "portfolio not found"})
	}

	// curl -i http://localhost:8080/api/v1/portfolios/1
}

func PostPortfolio(c *gin.Context) {
	var portfolio Portfolio
	//x, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Printf("%s", string(x))

	c.Bind(&portfolio)

	log.Println(portfolio)

	if portfolio.Created_by == "" {
		portfolio.Created_by = "Anonymous"
	}

	if portfolio.Updated_by == "" {
		portfolio.Updated_by = portfolio.Created_by
	}

	if portfolio.Author == "" {
		portfolio.Author = portfolio.Created_by
	}
	log.Println(portfolio.Name)
	if portfolio.Created_by != "" && portfolio.Name != "" && portfolio.Text_html != "" && portfolio.Updated_by != "" {

		if insert, _ := dbmap.Exec(`INSERT INTO Portfolio (createdby, updatedby, name, description, texthtml, demourl, author) VALUES (?, ?, ?, ?, ?, ?, ?)`, portfolio.Created_by, portfolio.Updated_by, portfolio.Name, portfolio.Description, portfolio.Text_html, portfolio.Demo_url, portfolio.Author); insert != nil {
			portfolio_id, err := insert.LastInsertId()
			if err == nil {
				content := &Portfolio{
					Id:           portfolio_id,
					Created_time: portfolio.Created_time,
					Updated_time: portfolio.Updated_time,
					Created_by:   portfolio.Created_by,
					Updated_by:   portfolio.Updated_by,
					Name:         portfolio.Name,
					Description:  portfolio.Description,
					Text_html:    portfolio.Text_html,
					Demo_url:     portfolio.Demo_url,
					Author:       portfolio.Author,
				}
				c.JSON(201, content)
				return
			} else {
				checkErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(400, gin.H{"error": "Fields are empty"})
		return
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/portfolios
}

func UpdatePortfolio(c *gin.Context) {
	id := c.Params.ByName("id")
	var portfolio Portfolio
	err := dbmap.SelectOne(&portfolio, "SELECT * FROM Portfolio WHERE id=?", id)

	if err == nil {
		var json Portfolio
		c.Bind(&json)

		portfolio_id, _ := strconv.ParseInt(id, 0, 64)

		portfolio := Portfolio{
			Id:           portfolio_id,
			Created_time: json.Created_time,
			Updated_time: json.Updated_time,
			Created_by:   json.Created_by,
			Updated_by:   json.Updated_by,
			Name:         json.Name,
			Description:  json.Description,
			Text_html:    json.Text_html,
			Demo_url:     json.Demo_url,
			Author:       json.Author,
		}

		if portfolio.Created_by == "" {
			portfolio.Created_by = "Anonymous"
		}

		if portfolio.Updated_by == "" {
			portfolio.Updated_by = portfolio.Created_by
		}

		if portfolio.Author == "" {
			portfolio.Author = portfolio.Created_by
		}

		if portfolio.Created_by != "" && portfolio.Updated_by != "" && portfolio.Name != "" && portfolio.Text_html != "" {
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
	err := dbmap.SelectOne(&portfolio, "SELECT * FROM Portfolio WHERE id=?", id)

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
