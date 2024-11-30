package testdata

import (
	"go-conv/testdata/domain"
	"go-conv/testdata/model"
)

// :convergen
//
//go:generate go run github.com/reedom/convergen
type Convergen interface {
	// :conv fromDomainCategory Category
	// :conv urlsToStrings PhotoUrls
	//DomainToModel(*domain.Pet) *model.Pet
	// :conv toDomainCategory Category
	// :conv stringsToURLs PhotoUrls
	// :conv domain.NewPetStatusFromValue Status
	ModelToDomain(*model.Pet) *domain.Pet
}

// ModelToDomain
// go-conv:conv
// go-conv:apply ArrayToSlice:slice
var ModelToDomain func(*model.Pet) *domain.Pet

func urlsToStrings(list []domain.URL) []string {
	ret := make([]string, len(list))
	for i, url := range list {
		ret[i] = url.String()
	}
	return ret
}

func stringsToURLs(list []string) []domain.URL {
	ret := make([]domain.URL, len(list))
	for i, s := range list {
		ret[i] = domain.NewURL(s)
	}
	return ret
}
