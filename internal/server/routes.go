package server

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (server *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	router.GET("/health", func(client *gin.Context) {
		// O Gin transforma esse mapa em um JSON limpo de forma nativa
		client.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"message": "API de Hotelaria respondendo com sucesso!",
		})
	})

	// 3. Rota de teste: Valida se o Gin consegue alcançar o contêiner do Postgres
	router.GET("/db-check", func(client *gin.Context) {
		// se o usuário cancelar a requisição, o banco pare de processar imediatamente.
		err := server.db.Ping(client.Request.Context())
		if err != nil {
			client.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"detail": "Não foi possível alcançar o banco de dados Postgres no Docker",
			})
			return
		}

		client.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"detail": "PostgreSQL integrado e respondendo às requisições do Gin com sucesso!",
		})
	})

	return router
}