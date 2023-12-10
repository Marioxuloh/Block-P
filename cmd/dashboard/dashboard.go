package Dashboard

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed view/*
var content embed.FS

func Dashboard(address string) {
	router := gin.Default()

	//manejo de errores en el middelware
	router.Use(gin.Recovery())

	// Servir archivos estáticos incrustados desde la carpeta "view"
	router.StaticFS("/dashboard", http.FS(content))

	// Configurar rutas
	router.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run(address)
}
