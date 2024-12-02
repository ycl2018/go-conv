package testdata

import (
	"github.com/ycl2018/go-conv/testdata/domain"
	"github.com/ycl2018/go-conv/testdata/model"
)

// ModelToDomain
// go-conv:generate
// go-conv:apply ArrayToSlice:slice
var ModelToDomain func(*model.Pet) *domain.Pet
