package clientRouter

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	Nstore "hack/internal/pkg/store"
	"net/http"
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
	engine.GET("/get", r.GetNews)
}

func (r *Router) GetNews(c *gin.Context) {
	news, err := r.store.GetNews(c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(news)
	c.IndentedJSON(http.StatusOK, news)
}

func (r *Router) Run() {
	_ = r.ginContext.Run("localhost:8080")
}
