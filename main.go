package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	StopWords       []string
	DeleteCharacter []string
}

type Index struct {
	Name     string
	Doc      []string
	Map      map[string][]int // {word: [doc_id1, doc_id2, ...]}
	Analyzer func(text string) []string
	Config   Config
}

var indices = make(map[string]*Index)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/_cat/indices", func(c echo.Context) error {
		res := make([]string, 0)
		for k := range indices {
			res = append(res, k)
		}
		return c.JSON(http.StatusOK, res)
	})

	e.POST("/:IndexName", func(c echo.Context) error {
		indexName := c.Param("IndexName")
		fmt.Println("indexName", indexName)
		index := newIndex(indexName, Config{})
		indices[indexName] = index
		return c.String(http.StatusOK, "Index created")
	})
	e.DELETE("/:IndexName", func(c echo.Context) error {
		indexName := c.Param("IndexName")
		_, ok := indices[indexName]
		if !ok {
			return c.String(http.StatusNotFound, "Index not found")
		}
		delete(indices, indexName)
		return c.String(http.StatusOK, "Index deleted")
	})

	e.POST("/:IndexName/_doc", func(c echo.Context) error {
		indexName := c.Param("IndexName")
		index, ok := indices[indexName]
		if !ok {
			return c.String(http.StatusNotFound, "Index not found")
		}
		// bodyを読み取る
		body := c.Request().Body
		defer body.Close()
		b, err := io.ReadAll(body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
		// jsonをmapに変換
		var data map[string]string
		err = json.Unmarshal(b, &data)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
		text, ok := data["text"]
		if !ok {
			return c.String(http.StatusBadRequest, "Text is required")
		}
		index.addIndex(text)
		return c.String(http.StatusOK, "Index added")
	})

	e.GET("/:IndexName", func(c echo.Context) error {
		indexName := c.Param("IndexName")
		index, ok := indices[indexName]
		if !ok {
			return c.String(http.StatusNotFound, "Index not found")
		}
		res, err := index.getIndex()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
		return c.JSON(http.StatusOK, res)
	})
	e.GET("/:IndexName/_search", func(c echo.Context) error {
		indexName := c.Param("IndexName")
		index, ok := indices[indexName]
		if !ok {
			return c.String(http.StatusNotFound, "Index not found")
		}
		body := c.Request().Body
		defer body.Close()
		b, err := io.ReadAll(body)
		if err != nil {
			e.Logger.Error(err)
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
		var data map[string]string
		fmt.Println("b", string(b))
		err = json.Unmarshal(b, &data)
		if err != nil {
			e.Logger.Error(err)
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
		word, ok := data["word"]
		if !ok {
			return c.String(http.StatusBadRequest, "Word is required")
		}
		res := index.searchIndex(word)
		return c.JSON(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func newIndex(name string, config Config) *Index {
	defaultConfig := Config{
		StopWords: []string{
			"a",
			"an",
			"the",
		},
		DeleteCharacter: []string{
			"!",
			",",
		},
	}
	if config.StopWords == nil {
		config.StopWords = defaultConfig.StopWords
	}
	if config.DeleteCharacter == nil {
		config.DeleteCharacter = defaultConfig.DeleteCharacter
	}
	return &Index{
		Name:   name,
		Doc:    make([]string, 0),
		Map:    make(map[string][]int),
		Config: config,
	}
}

func (i *Index) getIndex() (string, error) {
	bytes, err := json.Marshal(i.Map)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (i *Index) analyzeText(text string) []string {
	for _, word := range i.Config.DeleteCharacter {
		text = strings.ReplaceAll(text, word, "")
	}
	words := strings.Fields(text)
	for _, word := range i.Config.StopWords {
		for i := 0; i < len(words); i++ {
			if words[i] == word {
				words = append(words[:i], words[i+1:]...)
				i--
			}
		}
	}
	return words
}

func (i *Index) addIndex(text string) {
	i.Doc = append(i.Doc, text)
	index := len(i.Doc) - 1
	words := i.analyzeText(text)
	// split the text into words
	for _, word := range words {
		// check if the word is already in the map
		if _, ok := i.Map[word]; !ok {
			i.Map[word] = make([]int, 0)
		}
		i.Map[word] = append(i.Map[word], index)
	}
}

func unique(arr []int) []int {
	m := make(map[int]bool)
	uniqueArr := make([]int, 0)
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v] = true
			uniqueArr = append(uniqueArr, v)
		}
	}
	return uniqueArr
}

func (i *Index) searchIndex(text string) []string {
	words := i.analyzeText(text)
	var docIds []int
	for _, word := range words {
		docIds = append(docIds, i.Map[word]...)
	}
	if len(docIds) == 0 {
		return []string{}
	}
	docIds = unique(docIds)
	docs := make([]string, 0)
	for _, id := range docIds {
		docs = append(docs, i.Doc[id])
	}
	return docs
}
