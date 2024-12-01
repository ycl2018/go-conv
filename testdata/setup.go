package testdata

import (
	"go-conv/testdata/domain"
	"go-conv/testdata/model"
)

// ModelToDomain
// go-conv:generate
// go-conv:apply ArrayToSlice:slice
var ModelToDomain func(*model.Pet) *domain.Pet
