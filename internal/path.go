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

func (p *path) pop() {
	*p = (*p)[:len(*p)-1]
}

func (p *path) matchIgnore(ignoreType IgnoreType) bool {
	i, j := len(*p)-1, len(ignoreType.fields)-1
	if i < j {
		return false
	}
	for j >= 0 {
		if (*p)[i].name != ignoreType.fields[j] {
			return false
		}
		j--
		i--
	}
	// check type
	if ignoreType.typ != (*p)[i+1].structName {
		return false
	}
	return true
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
		for j >= 0 {
			l := (*p)[i].name
			r := ss[j]
			if l != r {
				return false
			}
			j--
			i--
		}
		return true
	}
	return false
}

func (p *path) matchFilter(filter Filter, src *types.Var) bool {
	if src.Type().String() != filter.typ {
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
		for j >= 0 {
			if (*p)[i].name != ss[j] {
				continue
			}
			j--
			i--
		}
		return true
	}
	return false
}
