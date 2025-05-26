package internal

import (
	"fmt"
	"go/types"
	"strings"
)

type fieldStep struct {
	src, dst field
}

type field struct {
	name       string
	structName string
}

type path []fieldStep

func (p *path) Top() *fieldStep {
	if len(*p) == 0 {
		return nil
	}
	return &(*p)[len(*p)-1]
}

func (p *path) Push(s fieldStep) {
	*p = append(*p, s)
}

func (p *path) Pop() {
	*p = (*p)[:len(*p)-1]
}

func (p *path) matchIgnore(ignoreType IgnoreType, src, dst types.Type) bool {
	if len(ignoreType.Fields) == 0 {
		switch ignoreType.IgnoreSide {
		case SideSrc:
			return p.matchFilter(Filter{
				Typ:   ignoreType.Tye,
				Paths: ignoreType.Paths,
			}, src, SideSrc)
		case SideDst:
			return p.matchFilter(Filter{
				Typ:   ignoreType.Tye,
				Paths: ignoreType.Paths,
			}, dst, SideDst)
		default:
			panic(fmt.Sprintf("unknown ignore side:%d", ignoreType.IgnoreSide))
		}
	}
	i := len(*p) - 1
	if i < 0 {
		return false
	}

	var matchField field
	switch side := ignoreType.IgnoreSide; side {
	case SideSrc:
		matchField = (*p)[i].src
	case SideDst:
		matchField = (*p)[i].dst
	default:
		panic(fmt.Sprintf("unknown ignore side:%d", side))
	}
	// check type
	if ignoreType.Tye != matchField.structName {
		return false
	}
	// check field
	fieldMatch := false
	for _, fieldName := range ignoreType.Fields {
		if matchField.name == fieldName {
			fieldMatch = true
			break
		}
	}
	if !fieldMatch {
		return false
	}
	return p.matchPath(ignoreType.Paths, ignoreType.IgnoreSide)
}

func (p *path) matchTransfer(transfer Transfer, dst *types.Var, src *types.Var) bool {
	if src.Type().String() != transfer.From || dst.Type().String() != transfer.To {
		return false
	}
	return p.matchPath(transfer.Paths, SideSrc) // transfer only apply to src side
}

func (p *path) matchFilter(filter Filter, src types.Type, side Side) bool {
	elemType, _, _ := dePointer(src)
	// auto de-pointer
	if elemType.String() != filter.Typ {
		return false
	}
	return p.matchPath(filter.Paths, side)
}

func (p *path) matchPath(paths []string, side Side) bool {
	if len(paths) == 0 {
		return true
	}
	for _, expr := range paths {
		ss := strings.Split(expr, ".")
		i, j := len(*p)-1, len(ss)-1
		if i != j {
			continue
		}
		var valid = true
		for j >= 0 {
			var matchName string
			switch side {
			case SideSrc:
				matchName = (*p)[i].src.name
			case SideDst:
				matchName = (*p)[i].dst.name
			default:
				panic(fmt.Sprintf("unknown side:%d", side))
			}
			if matchName != ss[j] {
				valid = false
				break
			}
			j--
			i--
		}
		if valid {
			return true
		}
	}
	return false
}
