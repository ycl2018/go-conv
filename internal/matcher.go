package internal

type FieldMatcher struct {
	conf map[string]map[string]string
}

func NewFieldMatcher() *FieldMatcher {
	return &FieldMatcher{conf: map[string]map[string]string{}}
}

func (f *FieldMatcher) AddMatch(structType, from, to string) {
	match, ok := f.conf[structType]
	if !ok || match == nil {
		match = make(map[string]string)
	}
	if _, ok := match[from]; ok {
		DefaultLogger.Printf("WARN: rematch %s's field:%s to %s", structType, from, to)
	} else {
		DefaultLogger.Printf("match %s's field:%s to %s", structType, from, to)
	}
	match[from] = to
	f.conf[structType] = match
}

func (f *FieldMatcher) HasMatch(structType, field string) (string, bool) {
	match, ok := f.conf[structType]
	if !ok || match == nil {
		return "", false
	}
	if has, ok := match[field]; ok {
		return has, true
	}
	return "", false
}
