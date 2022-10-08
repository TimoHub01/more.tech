package parserRouter

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Parser interface {
	Parse(ctx context.Context)
}

type Router struct {
	ginContext *gin.Engine
	parser     []Parser
}

func NewRouter(p ...Parser) *Router {
	return &Router{ginContext: gin.Default(), parser: p}
}

func (r *Router) SetUpRouter(engine *gin.Engine) {
	engine.POST("/parse", r.parseNews)
}

func (r *Router) parseNews(c *gin.Context) {
	for _, p := range r.parser {
		p.Parse(c)
	}
}

func (r *Router) parseNewsFromConsultant(c *gin.Context) {
	r.parser[1].Parse(c)
}

func (r *Router) Run() {
	_ = r.ginContext.Run("localhost:8080")
}
