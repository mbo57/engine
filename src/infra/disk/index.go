package disk

import (
	"app/domain/model"
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
	fmt.Println(err)
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
	// Write index to file
	// mk directory if not exist
	_, err := os.Stat(ii.rootPath)
	if os.IsNotExist(err) {
		os.Mkdir(ii.rootPath, 0755)
	}

	exist := false
	indexPath := ii.rootPath + "/" + i.Name

	var f *os.File
	_, err = os.Create(indexPath)
	if os.IsNotExist(err) {
		os.Mkdir(ii.rootPath, 0755)
		f, err = os.Create(indexPath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		exist = true
		f, err = os.Open(indexPath)
		if err != nil {
			return err
		}
	}

	if exist {
		header, err := readHeader(f)
		if err != nil {
			return err
		}
		if header.docCount > i.DocCount {
			i.DocCount = header.docCount
		}
		if header.textSize > i.TextSize {
			// TODO: rewrite
			return errors.New("元のtextSizeより大きいtextSizeは設定できません")
			// i.TextSize = header.textSize
		}
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
		err = binary.Write(f, binary.LittleEndian, textBytes)
		if err != nil {
			return err
		}

	}
	return nil
}

func readHeader(f *os.File) (header, error) {
	// Read header
	// header := make([]byte, HeaderSize)
	var docCount uint32
	if err := binary.Read(f, binary.LittleEndian, &docCount); err != nil {
		return header{}, err
	}
	var textSize uint32
	if err := binary.Write(f, binary.LittleEndian, &textSize); err != nil {
		return header{}, err
	}
	return header{docCount, textSize}, nil
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
