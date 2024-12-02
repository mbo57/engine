package valkey

import (
	"github.com/valkey-io/valkey-go"
)

type valkeyIndexInfra struct {
	client *valkey.Client
}

func NewValkeyIndexInfra(client *valkey.Client) *valkeyIndexInfra {
	return &valkeyIndexInfra{
		client: client,
	}
}

// func (v *valkeyIndexInfra) IndexWriter(indexName string, indexType string, indexFields []string) error {
// 	ctx := context.Background()
// 	v.client.Do(ctx, v.client.B())
// 	index := valkey.Index{
// 		Name:   indexName,
// 		Type:   indexType,
// 		Fields: indexFields,
// 	}
//
// 	err := v.client.CreateIndex(index)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
