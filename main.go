package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	//var url string
	r := gin.Default()

	r.GET("/ping/:ci", func(c *gin.Context) {
		var patterns = [9]string{"Cédula:", "Nombre:", "Estado:", "Municipio:", "Parroquia:", "Centro:", "Dirección:", "SERVICIO ELECTORAL", "VOTA en la Elección Constituyente"}
		var url = html.UnescapeString("http://www.cne.gov.ve/web/registro_electoral/ce.php?nacionalidad=:nacionalidad&cedula=:cedula")
		var contentFormatted = []string{}
		/*seeting values for url*/
		url = strings.Replace(url, ":nacionalidad", "VE", 1)
		url = strings.Replace(url, ":cedula", c.Param("ci"), 1)

		/*Get requesto to url*/
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		/*Get info in HTML*/
		dataInBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		pageContent := string(dataInBytes)

		/*Clear Tags Html*/
		pageContent = stripTags(pageContent)
		pageContent = strings.Replace(pageContent, "\t", "", -1)
		pageContent = strings.Replace(pageContent, "\n", "", -1)

		/*Find if exist Register*/
		titleExist := "ADVERTENCIA"
		registered := true
		if strings.Index(pageContent, titleExist) >= 0 {
			registered = false
		}

		/*if user registered*/
		if registered {
			/* split headers for separate with "|" */
			for i := 0; i < len(patterns); i++ {
				pageContent = strings.Replace(pageContent, patterns[i], "|", -1)
			}

			/*Convert string to Array*/
			contentFormatted = strings.Split(pageContent, "|")

			c.JSON(200, gin.H{
				"cedula":       c.Param("ci"),
				"inscrito":     true,
				"nacionalidad": "VE",
				"nombres":      contentFormatted[2],
				"estado":       contentFormatted[3],
				"municipio":    contentFormatted[4],
				"parroquia":    contentFormatted[5],
				"direccion":    contentFormatted[7],
			})
		} else {
			c.JSON(500, gin.H{
				"cedula":       c.Param("ci"),
				"nacionalidad": "VE",
				"inscrito":     false,
			})
		}

	})
	r.Run()
}

func stripTags(content string) string {
	re := regexp.MustCompile(`<(.|\n)*?>`)
	return re.ReplaceAllString(content, "")
}
