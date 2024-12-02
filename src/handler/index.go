package handler

import (
	"app/usecase"
	"app/util/logger"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
	logger  logger.Log
	usecase usecase.IndexUsecase
}

func NewIndexHandler(logger logger.Log, usecase usecase.IndexUsecase) *IndexHandler {
	return &IndexHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *IndexHandler) CatIndices(c echo.Context) error {
	indices, err := h.usecase.CatIndices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string][]string{"indices": indices})
}

func (h *IndexHandler) CreateIndex(c echo.Context) error {
	indexName := c.Param("IndexName")
	err := h.usecase.CreateIndex(indexName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "index created"})
}

func (h *IndexHandler) DeleteIndex(c echo.Context) error {
	indexName := c.Param("IndexName")
	err := h.usecase.DeleteIndex(indexName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "index deleted"})
}

func (h *IndexHandler) IndexDoc(c echo.Context) error {
	indexName := c.Param("IndexName")
	doc := make(map[string]interface{})
	if err := json.NewDecoder(c.Request().Body).Decode(&doc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	err := h.usecase.IndexAddDoc(indexName, doc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "index add doc"})
}

func (h *IndexHandler) GetIndexInfo(c echo.Context) error {
	indexName := c.Param("IndexName")
	indexDocs, indexMap, err := h.usecase.GetIndexInfo(indexName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"docs": indexDocs, "map": indexMap})
}

func (h *IndexHandler) SearchIndex(c echo.Context) error {
	indexName := c.Param("IndexName")
	query := make(map[string]interface{})
	if err := json.NewDecoder(c.Request().Body).Decode(&query); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	result, err := h.usecase.SearchIndex(indexName, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"result": result})
}

func (h *IndexHandler) CommitIndex(c echo.Context) error {
	indexName := c.Param("IndexName")
	err := h.usecase.CommitIndex(indexName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "commit success"})
}
