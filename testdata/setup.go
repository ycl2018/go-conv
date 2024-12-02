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
// go-conv:conv
//var ModelToDomain func(*model.PetNew) *domain.PetNew
