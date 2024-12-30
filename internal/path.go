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
	for _, field := range ignoreType.Fields {
		if (*p)[i].name == field {
			return true
		}
	}
	return false
}

func (p *path) matchTransfer(transfer Transfer, dst *types.Var, src *types.Var) bool {
	if src.Type().String() != transfer.From || dst.Type().String() != transfer.To {
		return false
	}
	if len(transfer.Paths) == 0 {
		return true
	}
	for _, expr := range transfer.Paths {
		ss := strings.Split(expr, ".")
		i, j := len(*p)-1, len(ss)-1
		if i != j {
			continue
		}
		var valid = true
		for j >= 0 {
			l := (*p)[i].name
			r := ss[j]
			if l != r {
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

func (p *path) matchFilter(filter Filter, src types.Type) bool {
	if src.String() != filter.Typ {
		return false
	}
	if len(filter.Paths) == 0 {
		return true
	}
	for _, expr := range filter.Paths {
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
