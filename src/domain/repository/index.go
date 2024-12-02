package repository

import "app/domain/model"

type IndexRepository interface {
	IndexWriter(index *model.Index) error
	GetIndexDocCount(indexName string) (uint32, error)
	GetDocs(indexName string, docIds []uint32) ([]model.Doc, error)
	DeleteDocs(indexName string, docIds []uint32) ([]uint32, error)
}
