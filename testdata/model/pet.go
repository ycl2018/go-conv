package model

import (
	"github.com/ycl2018/go-conv/testdata/common"
	"github.com/ycl2018/go-conv/testdata/model/domain"
)

type Pet struct {
	ID uint `storage:"id"`
	//Category  Category `storage:"category"`
	Name    string  `storage:"name"`
	NamePtr *string `storage:"namePtr"`
	//PhotoUrls []string `storage:"photoUrls"`
	Status             *string             `storage:"status"`
	Array              [3]*domain.Category `storage:"children"`
	Slices             []*domain.Category  `storage:"children"`
	Maps               map[string]*domain.Category
	Next               *Pet
	PtrToStruct        *domain.Category
	StructToPtr        domain.Category
	SlicesStruct       []domain.Category   `storage:"children"`
	SlicesPtr          []*domain.Category  `storage:"children"`
	ArrayToSlice       []*domain.Category  `storage:"children"`
	SliceToArray       [3]*domain.Category `storage:"children"`
	UnSupported        string
	StringConvert      MyString
	StringConvert2     string
	NumberCast         int
	ByteSliceToString  []byte
	ByteSliceToString2 MyString
	MapStringString    map[string]string
	Embed
	C       string
	D       int
	Common  *common.Common
	Common2 *common.Common
	Common3 common.Common
}

type PetNew struct {
	Common *common.Common
}

type MyString string

type Embed struct {
	A string
	B int
}
