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

func fromDomainCategory(cat domain.Category) model.Category {
	return model.Category{
		CategoryID: uint64(cat.ID),
		Name:       cat.Name,
	}
}

func toDomainCategory(cat model.Category) domain.Category {
	return domain.Category{
		ID:   uint(cat.CategoryID),
		Name: cat.Name,
	}
}

func DomainToModel(src *domain.Pet) (dest *model.Pet) {
	if src != nil {
		dest = new(model.Pet)
		dest.ID = src.ID
		dest.Name = src.Name
		dest.Status = string(src.Status)
	}
	return
}

func ModelToDomain(src *model.Pet) (dest *domain.Pet) {
	if src != nil {
		dest = new(domain.Pet)
		dest.ID = src.ID
		dest.Name = src.Name
		dest.Status = domain.PetStatus(src.Status)
	}
	return
}
