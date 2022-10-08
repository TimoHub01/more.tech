package parserRouter

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

type Parser interface {
	Parse(ctx context.Context, days int)
}

type Router struct {
	ginContext *gin.Engine
	ticker     *time.Ticker
	parser     []Parser
}

func NewRouter(p ...Parser) *Router {
	return &Router{ginContext: gin.Default(), parser: p, ticker: time.NewTicker(24 * time.Hour)}
}

func (r *Router) SetUpRouter(engine *gin.Engine) {
	engine.POST("/parse", r.parseNews)
}

func (r *Router) parseNews(c *gin.Context) {
	go r.RunUpdateService(c)
	for _, p := range r.parser {
		p.Parse(c, 365)
	}
}

func (r *Router) Run() {
	_ = r.ginContext.Run("localhost:8080")
}

func (r *Router) RunUpdateService(c *gin.Context) {
	for range r.ticker.C {
		for _, p := range r.parser {
			p.Parse(c, 1)
		}
	}
}
