package router

import (
	"app/handler"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, h *handler.IndexHandler) {
	e.GET("/_cat/indices", h.CatIndices)
	e.POST("/:IndexName/_doc", h.IndexDoc)
	e.POST("/:IndexName", h.CreateIndex)
	e.DELETE("/:IndexName", h.DeleteIndex)
	e.GET("/:IndexName", h.GetIndexInfo)
	e.GET("/:IndexName/_doc/:DocID", h.GetDoc)
	e.DELETE("/:IndexName/_doc/:DocID", h.DeleteDoc)
	e.GET("/:IndexName/_search", h.SearchIndex)
	e.POST("/:IndexName/_commit", h.CommitIndex)
}
