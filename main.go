package queryBuilder

import (
	"encoding/json"
	"strings"
)

type builder struct {
	query  any
	source []string
	sort   []map[string]any
	size   int
	from   int
}

func New() *builder {
	return &builder{}
}

type Generatable interface {
	generate() any
}

type Sort struct {
	Field string
	Order string // asc or desc
}

func (b *builder) Build() (string, error) {
	body := struct {
		Size   int              `json:"size,omitempty"`
		From   int              `json:"from,omitempty"`
		Query  any              `json:"query,omitempty"`
		Sort   []map[string]any `json:"sort,omitempty"`
		Source []string         `json:"_source,omitempty"`
	}{
		b.size,
		b.from,
		b.query,
		b.sort,
		b.source,
	}

	query, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(query), nil
}

func (b *builder) Size(value int) *builder {
	b.size = value
	return b
}

func (b *builder) From(value int) *builder {
	b.from = value
	return b
}

func (b *builder) Source(value []string) *builder {
	b.source = value
	return b
}

func (b *builder) Query(query Generatable) *builder {
	b.query = query.generate()
	return b
}

func (b *builder) Sort(sort ...Sort) *builder {
	if sort != nil && len(sort) > 0 {
		sortList := make([]map[string]any, len(sort))
		for i, s := range sort {
			order := map[string]string{"order": s.Order}
			m := map[string]any{}
			m[s.Field] = order
			sortList[i] = m
		}
		b.sort = append(b.sort, sortList...)
	}

	return b
}

type functionScore struct {
	query     any
	functions any
}

func (f *functionScore) generate() any {
	return struct {
		FunctionScore any `json:"function_score"`
	}{
		struct {
			Query     any `json:"query,omitempty"`
			Functions any `json:"functions,omitempty"`
		}{
			f.query,
			f.functions,
		},
	}
}

type Function struct {
	Filter Generatable
	Weight float32
}

func FunctionScore(query Generatable, functions []Function) Generatable {
	type functionType struct {
		Filter any     `json:"filter"`
		Weight float32 `json:"weight"`
	}
	_functions := make([]functionType, len(functions))

	for i, f := range functions {
		_functions[i] = functionType{f.Filter.generate(), f.Weight}
	}

	return &functionScore{
		query.generate(),
		_functions,
	}
}

type matchAll struct {
}

func (m *matchAll) generate() any {
	return struct {
		MatchAll any `json:"match_all"`
	}{struct{}{}}
}

func MatchAll() Generatable {
	return &matchAll{}
}

type match struct {
	match map[string]string
}

func (t *match) generate() any {
	return struct {
		Match map[string]string `json:"match,omitempty"`
	}{t.match}
}

func Match(field string, value string) Generatable {
	m := map[string]string{}
	m[field] = value
	return &match{m}
}

type matchPhrase struct {
	matchPhrase map[string]string
}

func (m *matchPhrase) generate() any {
	return struct {
		MatchPhrase map[string]string `json:"match_phrase,omitempty"`
	}{m.matchPhrase}
}

func MatchPhrase(field string, value []string) Generatable {
	m := map[string]string{}
	m[field] = strings.Join(value, " ")
	return &matchPhrase{m}
}

type term struct {
	term map[string]any
}

func (t *term) generate() any {
	return struct {
		Term map[string]any `json:"term,omitempty"`
	}{t.term}
}

func Term(field string, value any) Generatable {
	m := map[string]any{}
	m[field] = value
	return &term{m}
}

type terms struct {
	terms map[string]any
}

func (t *terms) generate() any {
	return struct {
		Terms map[string]any `json:"terms,omitempty"`
	}{t.terms}
}

func Terms[T any](field string, values []T) Generatable {
	m := map[string]any{}
	m[field] = values
	return &terms{m}
}

type prefix struct {
	prefix map[string]string
}

func (t *prefix) generate() any {
	return struct {
		Prefix map[string]string `json:"prefix,omitempty"`
	}{t.prefix}
}

func Prefix(field string, values string) Generatable {
	m := map[string]string{}
	m[field] = values
	return &prefix{m}
}

type exists struct {
	fieldName string
}

func (e *exists) generate() any {
	return struct {
		Exists any `json:"exists"`
	}{
		struct {
			Field string `json:"field"`
		}{
			e.fieldName,
		},
	}
}

func Exists(field string) Generatable {
	return &exists{field}
}
