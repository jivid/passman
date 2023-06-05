package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jivid/passman/passman/passman"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"flag"
)

var (
	p *passman.Passman
	passmanDir string
)

type rawPassword struct {
	Site     string
	Username string
	Password string
}

func getAllPasswords(c *gin.Context) {
	passwords, _ := p.GetAllContents()
	c.IndentedJSON(http.StatusOK, passwords)
}

func getPassword(c *gin.Context) {
	password, _ := p.GetForSite(c.Param("site"))
	if password == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No password for site %s", c.Param("site"))})
	} else {
		c.IndentedJSON(http.StatusOK, password)
	}
}

func createPassword(c *gin.Context) {
	var password rawPassword
	if err := c.BindJSON(&password); err != nil {
		return
	}
	if err := p.Create(password.Site, password.Username, password.Password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not create password" + err.Error()})
	} else {
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Created password"})
	}
}

func main() {
	flag.StringVar(&passmanDir, "dir", ".", "Directory where passman should init")
	flag.Parse()
	p = passman.NewPassman(passmanDir)
	p.InitOrLoad()
	router := gin.Default()
	router.GET("/passwords", getAllPasswords)
	router.GET("/passwords/:site", getPassword)
	router.POST("/passwords", createPassword)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Run("0.0.0.0:8080")
}
