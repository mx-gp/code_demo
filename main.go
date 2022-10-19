package main

import (
	"database/sql"
	"time"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306/cetec")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return db
}

func main() {

	// initialise gin
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Enable cors
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-type", "X-Auth-Token", "Access-Control-Allow-Origin"},
		MaxAge:       60 * time.Second,
	}))

	// API Routes
	var r *gin.Engine
	personRoutes := r.Group("/person")
	{
		personRoutes.POST("/create", PersonPOST)
		personRoutes.GET("/:id/info", PersonGET)
	}

}

func PersonPOST(c *gin.Context) {
	db := Connect()
	defer db.Close()
	var request struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		City        string `json:"city"`
		State       string `json:"state"`
		Street1     string `json:"street1"`
		Street2     string `json:"street2"`
		ZipCode     string `json:"zip_code"`
	}

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"error":   err,
			"message": err.Error(),
			"success": false,
		})
		return
	}

	/* "name": "YOURNAME",
	"phone_number": "123-456-7890",
	"city" : "Sacramento",
	"state" : "CA",
	"street1": "112 Main St",
	"street2": "Apt 12",
	"zip_code": "12345", */

	result, err := db.Exec("INSERT INTO person(name) VALUES(?)", request.Name)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"error":   err,
			"message": err.Error(),
			"success": false,
		})
		return
	}

	PersonID, err := result.LastInsertId()
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"error":   err,
			"message": err.Error(),
			"success": false,
		})
		return
	}

	_, err = db.Exec("INSERT INTO phone(number, person_id) VALUES(?, ?)", request.PhoneNumber, PersonID)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"error":   err,
			"message": err.Error(),
			"success": false,
		})
		return
	}

	rslt, err := db.Exec("INSERT INTO address(city, state, street1, street2, zip_code) VALUES(?, ?, ?, ?, ?)", request.City, request.State, request.Street1, request.Street2, request.ZipCode)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"error":   err,
			"message": err.Error(),
			"success": false,
		})
		return
	}

	AddressID, err := rslt.LastInsertId()
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"error":   err,
			"message": err.Error(),
			"success": false,
		})
		return
	}

	_, err = db.Exec("INSERT INTO address_join(person_id, address_id) VALUES(?, ?)", PersonID, AddressID)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"error":   err,
			"message": err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"data":    "Data Inserted",
		"success": true,
	})

}

func PersonGET(c *gin.Context) {
	type PersonDetails struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		City        string `json:"city"`
		State       string `json:"state"`
		Street1     string `json:"street1"`
		Street2     string `json:"street2"`
		ZipCode     string `json:"zip_code"`
	}

	var p PersonDetails
	var Arr []PersonDetails

	db := Connect()
	defer db.Close()

	rows, err := db.Query("SELECT person.name, phone.number, address.city, address.state, address.street1, address.street2, address.zip_code from person join phone on person.id=phone.person_id join address_join on person.id=address_join.person_id join address on address_join.address_id=address.id")
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"error":   err,
			"message": err.Error(),
			"success": false,
		})
		return
	}

	for rows.Next() {
		err = rows.Scan(&p.Name, &p.PhoneNumber, &p.City, &p.State, &p.Street1, &p.Street2, &p.ZipCode)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"error":   err,
				"message": err.Error(),
				"success": false,
			})
			return
		} else {
			Arr = append(Arr, p)
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"data":    Arr,
		"success": true,
	})
}
