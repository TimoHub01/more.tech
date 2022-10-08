package clientRouter

import (
	"context"
	"github.com/gin-gonic/gin"
	Nstore "hack/internal/pkg/store"
)

type store interface {
	GetNews(ctx context.Context) ([]Nstore.New, error)
}

type Router struct {
	ginContext *gin.Engine
	store      store
}

func NewRouter(s store) *Router {
	return &Router{ginContext: gin.Default(), store: s}
}

func (r *Router) SetUpRouter(engine *gin.Engine) {
	engine.POST("/get", r.GetNews)
}

func (r *Router) GetNews(c *gin.Context) {
	r.store.GetNews(c)
}

func (r *Router) Run() {
	_ = r.ginContext.Run("localhost:8080")
}
