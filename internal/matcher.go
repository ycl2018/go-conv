package internal

type FieldMatcher struct {
	Conf map[string]map[string]string
}

func NewFieldMatcher() *FieldMatcher {
	return &FieldMatcher{Conf: map[string]map[string]string{}}
}

func (f *FieldMatcher) AddMatch(structType, from, to string) {
	match, ok := f.Conf[structType]
	if !ok || match == nil {
		match = make(map[string]string)
	}
	if _, ok := match[from]; ok {
		DefaultLogger.Printf("WARN: rematch %s's field:%s to %s", structType, from, to)
	}
	match[from] = to
	f.Conf[structType] = match
}

func (f *FieldMatcher) HasMatch(structType, field string) (string, bool) {
	match, ok := f.Conf[structType]
	if !ok || match == nil {
		return "", false
	}
	if has, ok := match[field]; ok {
		return has, true
	}
	return "", false
}
