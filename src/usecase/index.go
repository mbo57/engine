package usecase

import (
	"app/domain/model"
	"app/domain/repository"
	"app/util"
	"app/util/logger"
	"app/util/pserror"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type indexUsecase struct {
	logger    logger.Log
	indexRepo repository.IndexRepository
	indices   map[string]*model.Index
}

func NewIndexUsecase(
	logger logger.Log,
	indexRepo repository.IndexRepository,
) IndexUsecase {
	return &indexUsecase{
		logger:    logger,
		indexRepo: indexRepo,
		indices:   make(map[string]*model.Index),
	}
}

type IndexUsecase interface {
	IndexWriter() error
	CatIndices() ([]string, error)
	CreateIndex(indexName string) error
	DeleteIndex(indexName string) error
	IndexAddDoc(indexName string, doc map[string]interface{}) error
	GetDoc(indexName string, docID string) (*model.Doc, error)
	DeleteDoc(indexName string, docID string) error
	GetIndexInfo(indexName string) ([]model.Doc, map[string]map[uint32][]uint32, error)
	SearchIndex(indexName string, query map[string]interface{}) ([]result, error)
	CommitIndex(indexName string) error
}

func (i *indexUsecase) CatIndices() ([]string, error) {
	i.logger.Debug("CatIndices")
	indices := []string{}
	for name := range i.indices {
		indices = append(indices, name)
	}
	return indices, nil
}

func (i *indexUsecase) IndexWriter() error {
	i.logger.Info("IndexWriter")
	return nil
}

func (i *indexUsecase) CreateIndex(indexName string) error {
	i.logger.Debug("CreateIndex")
	if _, ok := i.indices[indexName]; ok {
		return errors.New("index already exists")
	}
	index := model.Index{
		Name:     indexName,
		Docs:     []model.Doc{},
		Map:      map[string]map[uint32][]uint32{},
		Analyzer: DefaultAnalyzer(DefualtStopWords, DefualtDeleteCharacter),
		DocCount: 0,
		TextSize: 2048,
	}
	i.indices[indexName] = &index
	return nil
}

func (i *indexUsecase) DeleteIndex(indexName string) error {
	i.logger.Debug("DeleteIndex")
	if _, ok := i.indices[indexName]; !ok {
		return errors.New("index not found")
	}
	delete(i.indices, indexName)
	return nil
}

func (i *indexUsecase) IndexAddDoc(
	indexName string,
	inputDoc map[string]interface{},
) error {
	i.logger.Debug("IndexDoc")
	index, ok := i.indices[indexName]
	if !ok {
		return errors.New("index not found")
	}
	inputText, ok := inputDoc["text"].(string)
	if !ok {
		return errors.New("text not found")
	}

	docCount, err := i.indexRepo.GetIndexDocCount(index.Name)
	if err != nil {
		i.logger.Error(err.Error())
		return err
	}
	if index.DocCount == 0 {
		index.DocCount = docCount
	}

	index.DocCount++
	id := index.DocCount - 1
	fmt.Println("new doc id", id)
	doc := model.NewDoc(id, inputText)
	index.Docs = append(index.Docs, doc)
	byteLen := len([]byte(inputText))
	if uint32(byteLen) > index.TextSize {
		index.TextSize = uint32(byteLen)
	}
	words := index.Analyzer(inputText)
	for i, word := range words {
		if _, ok := index.Map[word]; !ok {
			index.Map[word] = map[uint32][]uint32{}
		}
		if _, ok := index.Map[word][id]; !ok {
			index.Map[word][id] = []uint32{}
		}
		index.Map[word][id] = append(index.Map[word][id], uint32(i))
	}
	return nil
}

func (i *indexUsecase) GetIndexInfo(indexName string) ([]model.Doc, map[string]map[uint32][]uint32, error) {
	i.logger.Debug("GetIndexMap")

	index, ok := i.indices[indexName]
	if !ok {
		return []model.Doc{}, nil, errors.New("index not found")
	}
	return index.Docs, index.Map, nil
}

func (i *indexUsecase) GetDoc(
	indexName string,
	docID string,
) (*model.Doc, error) {
	i.logger.Debug("GetDoc")
	tmpID, err := strconv.ParseInt(docID, 10, 32)
	if err != nil {
		return nil, err
	}
	id := uint32(tmpID)
	index, okMemory := i.indices[indexName]
	if okMemory {
		// TODO: 高速化
		// Docsの構造を変える
		for _, doc := range index.Docs {
			if doc.Id == id {
				return &doc, nil
			}
		}
	}
	docs, err := i.indexRepo.GetDocs(indexName, []uint32{id})
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, pserror.NotFoundDoc
	}
	return &docs[0], nil
}

func (i *indexUsecase) DeleteDoc(
	indexName string,
	docID string,
) error {
	i.logger.Debug("DeleteDoc")
	tmpID, err := strconv.ParseInt(docID, 10, 32)
	if err != nil {
		return err
	}
	id := uint32(tmpID)
	index, okMemory := i.indices[indexName]
	if okMemory {
		// TODO: 高速化
		// Docsの構造を変える
		for _, doc := range index.Docs {
			if doc.Id == id {
				doc.Deleted = true
				return nil
			}
		}
	}

	_, err = i.indexRepo.DeleteDocs(indexName, []uint32{id})
	if err != nil {
		return err
	}
	return nil
}

func (i *indexUsecase) SearchIndex(
	indexName string,
	query map[string]interface{},
) ([]result, error) {
	i.logger.Debug("SearchIndex")
	index, ok := i.indices[indexName]
	if !ok {
		return nil, errors.New("index not found")
	}
	text, ok := query["text"].(string)
	if !ok {
		return nil, errors.New("text not found")
	}

	searchMode := "AND"
	if _, ok := query["search_mode"]; ok {
		searchMode, ok = query["search_mode"].(string)
		if !ok {
			return nil, errors.New("invalid search mode. please set AND or OR")
		}
		if searchMode != "AND" && searchMode != "OR" {
			return nil, errors.New("invalid search mode. please set AND or OR")
		}
	}

	words := index.Analyzer(text)
	macthWordIds := make([][]uint32, len(words))
	for i, word := range words {
		if _, ok := index.Map[word]; !ok {
			continue
		}
		for id := range index.Map[word] {
			macthWordIds[i] = append(macthWordIds[i], id)
		}
	}

	resultIds := []uint32{}
	if searchMode == "AND" {
		for i := 0; i < len(macthWordIds); i++ {
			sort.Slice(macthWordIds[i], func(j, k int) bool {
				return macthWordIds[i][j] < macthWordIds[i][k]
			})
		}
		resultIds = util.FindCommonElements(macthWordIds)
	} else {
		ids := []uint32{}
		for i := 0; i < len(macthWordIds); i++ {
			ids = append(ids, macthWordIds[i]...)
		}
		resultIds = util.FindUniqueElemens(ids)
	}

	var results []result
	for _, id := range resultIds {
		doc := index.Docs[id]
		if doc.Deleted {
			continue
		}
		results = append(results, result{Id: doc.Id, Text: doc.Text})
	}

	resultRank(&results)

	return results, nil
}

type result struct {
	Id   uint32 `json:"id"`
	Text string `json:"text"`
}

func resultRank(results *[]result) {
	sort.Slice(*results, func(i, j int) bool {
		return (*results)[i].Id < (*results)[j].Id
	})
}

func (i *indexUsecase) CommitIndex(indexName string) error {
	i.logger.Debug("CommitIndex")
	index, ok := i.indices[indexName]
	if !ok {
		return errors.New("index not found")
	}
	err := i.indexRepo.IndexWriter(index)
	if err != nil {
		i.logger.Error(err.Error())
		return err
	}
	delete(i.indices, indexName)
	return nil

}

var DefualtDeleteCharacter = []string{"!", "?", ".", ","}
var DefualtStopWords = []string{"a", "an", "the"}

func DefaultAnalyzer(
	stopWords []string,
	deleteCharacter []string,
) func(text string) []string {
	return func(text string) []string {
		text = strings.ToLower(text)
		for _, word := range deleteCharacter {
			text = strings.ReplaceAll(text, word, "")
		}
		words := strings.Fields(text)
		for _, word := range stopWords {
			for i := 0; i < len(words); i++ {
				if words[i] == word {
					words = append(words[:i], words[i+1:]...)
					i--
				}
			}
		}
		return words
	}
}
