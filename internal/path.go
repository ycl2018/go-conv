package internal

import (
	"go/types"
	"strings"
)

type fieldStep struct {
	name       string
	structName string
}

type path []fieldStep

func (p *path) Push(s fieldStep) {
	*p = append(*p, s)
}

func (p *path) Pop() {
	*p = (*p)[:len(*p)-1]
}

func (p *path) matchIgnore(ignoreType IgnoreType, srcType types.Type) bool {
	if len(ignoreType.Fields) == 0 {
		return p.matchFilter(Filter{
			Typ:   ignoreType.Tye,
			Paths: ignoreType.Paths,
		}, srcType)
	}
	i := len(*p) - 1
	if i < 0 {
		return false
	}
	// check type
	if ignoreType.Tye != (*p)[i].structName {
		return false
	}
	// check field
	fieldMatch := false
	for _, field := range ignoreType.Fields {
		if (*p)[i].name == field {
			fieldMatch = true
			break
		}
	}
	if !fieldMatch {
		return false
	}
	return p.matchPath(ignoreType.Paths)
}

func (p *path) matchTransfer(transfer Transfer, dst *types.Var, src *types.Var) bool {
	if src.Type().String() != transfer.From || dst.Type().String() != transfer.To {
		return false
	}
	return p.matchPath(transfer.Paths)
}

func (p *path) matchFilter(filter Filter, src types.Type) bool {
	if src.String() != filter.Typ {
		return false
	}
	return p.matchPath(filter.Paths)
}

func (p *path) matchPath(paths []string) bool {
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
			if (*p)[i].name != ss[j] {
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
