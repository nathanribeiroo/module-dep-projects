package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RouteMount encapsula a lógica de montagem de um conjunto de rotas em um router do Gin.
type RouteMount func(r gin.IRouter)

// Server é o ponto central de configuração e execução da API HTTP baseada em Gin.
type Server struct {
	gin         *gin.Engine
	ginMode     string
	middlewares []gin.HandlerFunc
	routes      []RouteMount
}

// N devolve uma instância limpa de Server pronta para ser configurada fluentemente.
func N() *Server {
	return &Server{
		gin:         nil,
		ginMode:     gin.ReleaseMode,
		middlewares: []gin.HandlerFunc{},
		routes:      []RouteMount{},
	}
}

// Middlewares registra middlewares globais que serão aplicados a todas as rotas.
func (s *Server) Middlewares(middleware ...gin.HandlerFunc) *Server {
	s.middlewares = append(s.middlewares, middleware...)
	return s
}

// Routes permite injetar funções de montagem de rotas no pipeline da aplicação.
func (s *Server) Routes(route ...RouteMount) *Server {
	s.routes = append(s.routes, route...)
	return s
}

// GinMode define o modo operacional do Gin (debug, release, test).
func (s *Server) GinMode(mode string) *Server {
	s.ginMode = mode
	return s
}

// Run inicializa o engine do Gin, aplica middlewares, monta rotas e expõe o servidor HTTP.
func (s *Server) Run(addr string) {
	gin.SetMode(s.ginMode)

	s.gin = gin.New()

	s.addInternalMiddlewares()
	s.gin.Use(s.middlewares...)

	s.addHealthCheck()

	for _, route := range s.routes {
		route(s.gin)
	}

	fmt.Println("HTTP server is running...")

	if err := s.gin.Run(":" + addr); err != nil {
		fmt.Errorf("Failed to start server: %v", err)
	}
}

// addHealthCheck registra o endpoint padrão de verificação de saúde da aplicação.
func (s *Server) addHealthCheck() {
	s.gin.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// addInternalMiddlewares aplica middlewares internos obrigatórios antes dos customizados.
func (s *Server) addInternalMiddlewares() {
	s.gin.Use(
		gin.Recovery(),
		addLogger(),
		xItauCorrelationId(),
	)
}
