package repository

import "app/domain/model"

type IndexRepository interface {
	IndexWriter(index *model.Index) error
	GetIndexDocCount(indexName string) (uint32, error)
}
