package test_gen

import (
	"go-conv/testdata/domain"
	"go-conv/testdata/model"
)

func ModelToDomain(src *model.Pet) (dst *domain.Pet) {
	if src != nil {
		dst = new(domain.Pet)
		dst.ID = src.ID
		dst.Name = src.Name
		if src.Status != nil {
			dst.Status = new(domain.PetStatus)
			*dst.Status = domain.PetStatus(*src.Status)
		}
		for i := 0; i < 3; i++ {
			if src.Children[i] != nil {
				dst.Children[i] = new(domain.Category)
				dst.Children[i].CategoryID = src.Children[i].CategoryID
				dst.Children[i].Name = src.Children[i].Name
				if src.Children[i].Foo != nil {
					dst.Children[i].Foo = new(domain.Foo)
					if src.Children[i].Foo.Bar != nil {
						dst.Children[i].Foo.Bar = new(string)
						*dst.Children[i].Foo.Bar = *src.Children[i].Foo.Bar
					}
				}
			}
		}
		if len(src.Childrens) > 0 {
			dst.Childrens = make([]*domain.Category, len(src.Childrens))
			for i := 0; i < len(src.Childrens); i++ {
				if src.Childrens[i] != nil {
					dst.Childrens[i] = new(domain.Category)
					dst.Childrens[i].CategoryID = src.Childrens[i].CategoryID
					dst.Childrens[i].Name = src.Childrens[i].Name
					if src.Childrens[i].Foo != nil {
						dst.Childrens[i].Foo = new(domain.Foo)
						if src.Childrens[i].Foo.Bar != nil {
							dst.Childrens[i].Foo.Bar = new(string)
							*dst.Childrens[i].Foo.Bar = *src.Childrens[i].Foo.Bar
						}
					}
				}
			}
		}
		if len(src.MapChild) > 0 {
			dst.MapChild = make(map[string]*domain.Category, len(src.MapChild))
			for k, v := range src.MapChild {
				var (
					tmpK string
					tmpV *domain.Category
				)
				tmpK = k
				if v != nil {
					tmpV = new(domain.Category)
					tmpV.CategoryID = v.CategoryID
					tmpV.Name = v.Name
					if v.Foo != nil {
						tmpV.Foo = new(domain.Foo)
						if v.Foo.Bar != nil {
							tmpV.Foo.Bar = new(string)
							*tmpV.Foo.Bar = *v.Foo.Bar
						}
					}
				}
				dst.MapChild[tmpK] = tmpV
			}
		}
	}
	return
}
