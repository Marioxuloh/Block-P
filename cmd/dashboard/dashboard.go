package Dashboard

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed View/*
var content embed.FS

func Dashboard(address string) {
	router := gin.Default()

	// Manejo de errores en el middleware
	router.Use(gin.Recovery())

	// Configurar ruta para servir archivos estáticos embebidos
	router.GET("/dashboard/:page?", func(c *gin.Context) {
		// Obtener el nombre de la página desde la URL
		page := c.Param("page")

		// Construir la ruta del archivo HTML embebido
		filePath := "View/" + page + ".html"

		// Leer el contenido del archivo HTML desde el sistema de archivos embebido
		fileContent, err := content.ReadFile(filePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error al leer el archivo "+page+".html")
			return
		}

		// Crear una plantilla HTML y renderizarla
		tmpl, err := template.New(page).Parse(string(fileContent))
		if err != nil {
			c.String(http.StatusInternalServerError, "Error al procesar la plantilla")
			return
		}

		// Puedes pasar datos adicionales a la plantilla si es necesario
		data := map[string]interface{}{
			// Datos adicionales si es necesario
		}

		// Renderizar la plantilla en la respuesta HTTP
		err = tmpl.Execute(c.Writer, data)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error al renderizar la plantilla")
			return
		}
	})

	router.Run(address)
}
