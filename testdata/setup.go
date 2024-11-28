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

func ModelToDomain(src *model.Pet) (dest *domain.Pet) {
	if src != nil {
		dest = new(domain.Pet)
		dest.ID = src.ID
		dest.Name = src.Name
		if src.Status != nil {
			dest.Status = new(domain.PetStatus)
			*dest.Status = domain.PetStatus(*src.Status)
		}
		for i := 0; i < 3; i++ {
			if src.Children[i] != nil {
				dest.Children[i] = new(domain.Category)
				dest.Children[i].CategoryID = src.Children[i].CategoryID
				dest.Children[i].Name = src.Children[i].Name
				if src.Children[i].Foo != nil {
					dest.Children[i].Foo = new(domain.Foo)
					if src.Children[i].Foo.Bar != nil {
						dest.Children[i].Foo.Bar = new(string)
						*dest.Children[i].Foo.Bar = *src.Children[i].Foo.Bar
					}
				}
			}
		}
	}
	return
}
