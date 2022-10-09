package clientRouter

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	Nstore "hack/internal/pkg/store"
	"net/http"
)

type store interface {
	GetNews(ctx context.Context, param int) ([]Nstore.New, error)
	GetTrends(ctx context.Context) ([]Nstore.Trend, error)
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
	engine.GET("/trends", r.GetTrends)
}

func (r *Router) GetNews(c *gin.Context) {
	var param int
	if err := c.BindJSON(&param); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	news, err := r.store.GetNews(c, param)
	if err != nil {
		fmt.Println(err)
	}
	c.IndentedJSON(http.StatusOK, news)
}

func (r *Router) GetTrends(c *gin.Context) {
	trends, err := r.store.GetTrends(c)
	if err != nil {
		fmt.Println(err)
	}
	c.IndentedJSON(http.StatusOK, trends)
}

func (r *Router) Run() {
	_ = r.ginContext.Run("localhost:8080")
}
