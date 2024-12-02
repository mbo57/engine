package disk

import (
	"app/domain/model"
	"app/util/pserror"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	typeDocs = iota
	typeMap
)

type header struct {
	docCount uint32
	textSize uint32
}

// Headerは以下の通り
// 4byte: doc数
// 4byte: textのbyte数
var HeaderSize uint32 = 8

type fileIndexInfra struct {
	rootPath string
}

func NewIndexRepository(rootPath string) *fileIndexInfra {
	return &fileIndexInfra{
		rootPath: rootPath,
	}
}

func (ii *fileIndexInfra) GetIndexDocCount(indexName string) (uint32, error) {
	indexPath := ii.rootPath + "/" + indexName
	f, err := os.Open(indexPath)
	if os.IsNotExist(err) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	header, err := readHeader(f)
	if err != nil {
		return 0, err
	}
	return header.docCount, nil
}

func (ii *fileIndexInfra) IndexWriter(i *model.Index) error {
	_, err := os.Stat(ii.rootPath)
	if os.IsNotExist(err) {
		os.Mkdir(ii.rootPath, 0755)
	}

	indexPath := ii.rootPath + "/" + i.Name
	f, err := os.OpenFile(indexPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	if stat.Size() != 0 {
		fmt.Println("exist", f.Name())
		header, err := readHeader(f)
		if err != nil {
			return err
		}
		if header.textSize > i.TextSize {
			// TODO: rewrite or create new file
			return errors.New("元のtextSizeより大きいtextSizeは設定できません")
			// i.TextSize = header.textSize
		}

	}

	// Write header
	header := header{
		docCount: i.DocCount,
		textSize: i.TextSize,
	}
	if err := writeHeader(f, header); err != nil {
		return err
	}

	dataSize := i.TextSize + 1
	for _, v := range i.Docs {
		textBytes := make([]byte, i.TextSize)
		copy(textBytes, []byte(v.Text))
		if v.Deleted {
			f.Seek(int64(HeaderSize+v.Id*dataSize), io.SeekStart)
		} else {
			f.Seek(0, io.SeekEnd)
		}
		if err := writeDoc(f, v, i.TextSize); err != nil {
			return err
		}

	}
	f.Close()
	return nil
}

func (ii *fileIndexInfra) GetDocs(indexName string, ids []uint32) ([]model.Doc, error) {
	indexPath := ii.rootPath + "/" + indexName
	f, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, pserror.NotFoundIndex
		}
		return nil, err
	}
	defer f.Close()

	header, err := readHeader(f)
	if err != nil {
		return nil, err
	}

	docs := []model.Doc{}
	dataSize := header.textSize + 1
	for _, id := range ids {
		f.Seek(int64(HeaderSize+id*dataSize), io.SeekStart)
		doc, err := readDoc(f, header.textSize)
		doc.Id = id
		if err != nil {
			return nil, err
		}
		// if doc.Deleted {
		// 	continue
		// }
		docs = append(docs, doc)
	}
	return docs, nil
}

func (ii *fileIndexInfra) DeleteDocs(indexName string, ids []uint32) ([]uint32, error) {
	indexPath := ii.rootPath + "/" + indexName
	f, err := os.OpenFile(indexPath, os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, pserror.NotFoundIndex
		}
		return nil, err
	}
	defer f.Close()

	header, err := readHeader(f)
	if err != nil {
		return nil, err
	}

	var deletedIds []uint32
	dataSize := header.textSize + 1
	for _, id := range ids {
		f.Seek(int64(HeaderSize+id*dataSize+header.textSize), io.SeekStart)
		if err := binary.Write(f, binary.LittleEndian, true); err != nil {
			return deletedIds, err
		}
		deletedIds = append(deletedIds, id)
	}
	if deletedIds == nil {
		return deletedIds, pserror.NotFoundDoc
	}

	return deletedIds, nil
}

func readHeader(f *os.File) (header, error) {
	// Read header
	// header := make([]byte, HeaderSize)
	f.Seek(0, io.SeekStart)
	var docCount uint32
	if err := binary.Read(f, binary.LittleEndian, &docCount); err != nil {
		return header{}, err
	}
	var textSize uint32
	if err := binary.Read(f, binary.LittleEndian, &textSize); err != nil {
		return header{}, err
	}
	return header{docCount, textSize}, nil
}

func writeHeader(f *os.File, header header) error {
	// Write header
	f.Seek(0, io.SeekStart)
	if err := binary.Write(f, binary.LittleEndian, header.docCount); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, header.textSize); err != nil {
		return err
	}
	return nil
}

func writeDoc(f *os.File, doc model.Doc, textSize uint32) error {
	// Write doc
	textBytes := make([]byte, textSize)
	copy(textBytes, []byte(doc.Text))
	if err := binary.Write(f, binary.LittleEndian, textBytes); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, doc.Deleted); err != nil {
		return err
	}
	return nil
}

func readDoc(f *os.File, textSize uint32) (model.Doc, error) {
	// Read doc
	textBytes := make([]byte, textSize)
	if err := binary.Read(f, binary.LittleEndian, textBytes); err != nil {
		return model.Doc{}, err
	}
	var deleted bool
	if err := binary.Read(f, binary.LittleEndian, &deleted); err != nil {
		return model.Doc{}, err
	}
	return model.Doc{
		Text:    string(bytes.TrimRight(textBytes, "\x00")),
		Deleted: deleted,
	}, nil
}
