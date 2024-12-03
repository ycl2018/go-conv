package testdata

import (
	"github.com/ycl2018/go-conv/testdata/domain"
	"github.com/ycl2018/go-conv/testdata/model"
)

// ModelToDomain
// go-conv:generate
// go-conv:conv
var ModelToDomain func(*model.Pet) *domain.Pet

// ModelToDomain
// go-conv:generate
// go-conv:copy
var ModelToDomain2 func(*model.Pet) *domain.Pet

// ModelToDomain
// go-conv:generate
// go-conv:conv
//var ModelToDomain func(*model.PetNew) *domain.PetNew

// go-conv:generate
// go-conv:copy
//var Foo2Bar func(src *Foo) *Bar

type Foo struct {
	Str     string
	Slice   []string
	Map     map[string]string
	Pointer string
	Alias   string
}

type Bar struct {
	Str     string
	Slice   []string
	Map     map[string]string
	Pointer *string
	Alias   StringAlias
}

type StringAlias string
