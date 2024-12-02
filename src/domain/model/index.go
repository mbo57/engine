package model

type Index struct {
	Name     string
	Docs     []Doc
	Map      map[string]map[uint32][]uint32
	Analyzer func(text string) []string
	DocCount uint32
	TextSize uint32
}

type Doc struct {
	Id      uint32
	Text    string
	Deleted bool
}

func NewDoc(id uint32, text string) Doc {
	return Doc{
		Id:      id,
		Text:    text,
		Deleted: false,
	}
}
