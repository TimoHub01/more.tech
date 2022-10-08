package app

import "github.com/gin-gonic/gin"

type Router struct {
	ginContext *gin.Engine
	subRouters []subRouter
}
type subRouter interface {
	SetUpRouter(engine *gin.Engine)
}

func NewRouter(subRouters ...subRouter) *Router {
	return &Router{gin.Default(), subRouters}
}

func (r *Router) SetUpRouter() {
	for _, s := range r.subRouters {
		s.SetUpRouter(r.ginContext)
	}
}

func (r *Router) Run() {
	_ = r.ginContext.Run("localhost:8080")
}
