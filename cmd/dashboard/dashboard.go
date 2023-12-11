package Dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
)

func Dashboard(address string) {
	r := gin.Default()

	// Cargar el sistema de archivos embebido
	pkger.Include("/view/build")

	// Servir archivos estáticos desde el sistema de archivos embebido
	r.StaticFS("/static", http.Dir(pkger.Dir("/view/build/static")))

	// Servir el archivo HTML principal
	r.NoRoute(func(c *gin.Context) {

		c.FileFromFS("/view/build/index.html", http.Dir(pkger.Dir("/view")))

	})

	// Manejar todas las demás solicitudes que no coinciden con un método HTTP específico
	r.NoMethod(func(c *gin.Context) {
		c.FileFromFS("/view/build/index.html", http.Dir(pkger.Dir("/view")))
	})

	r.Run(address)
}
