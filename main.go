package queryBuilder

import (
	"encoding/json"
	"strings"

	"github.com/aquasecurity/esquery"
)

type DataSource string

const (
	ES DataSource = "ElasticSearch"
)

type Builder struct {
	query       any
	source      []string
	sort        []map[string]any
	size        int
	from        int
	searchAfter []string
	aggs        map[string]map[string]any
}

func New() *Builder {
	return &Builder{}
}

type Generatable interface {
	generate() any
}

type Sort struct {
	Field string
	Order string // asc or desc
}

func (b *Builder) Build(dc DataSource) (string, error) {
	switch dc {
	case ES:
		fallthrough
	default:
		return b.buildElasticSearch()
	}
}

func (b *Builder) buildElasticSearch() (string, error) {
	body := struct {
		Size        int                       `json:"size,omitempty"`
		From        int                       `json:"from,omitempty"`
		Query       any                       `json:"query,omitempty"`
		Sort        []map[string]any          `json:"sort,omitempty"`
		Source      []string                  `json:"_source,omitempty"`
		SearchAfter []string                  `json:"search_after,omitempty"`
		Aggs        map[string]map[string]any `json:"aggs,omitempty"`
	}{
		b.size,
		b.from,
		b.query,
		b.sort,
		b.source,
		b.searchAfter,
		b.aggs,
	}

	query, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(query), nil
}

func (b *Builder) Size(value int) *Builder {
	b.size = value
	return b
}

func (b *Builder) From(value int) *Builder {
	b.from = value
	return b
}

func (b *Builder) Source(value []string) *Builder {
	b.source = value
	return b
}

func (b *Builder) Query(query Generatable) *Builder {
	b.query = query.generate()
	return b
}

func (b *Builder) Sort(sort ...Sort) *Builder {
	if len(sort) > 0 {
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

func (b *Builder) SearchAfter(values ...string) *Builder {
	b.searchAfter = values
	return b
}

func (b *Builder) Aggs(values ...Generatable) *Builder {
	aggs := make(map[string]map[string]any, len(values))
	for _, a := range values {
		for k, v := range a.generate().(map[string]map[string]any) {
			aggs[k] = v
		}
	}
	b.aggs = aggs
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

type rangeQuery struct {
	fieldName string
	params    RangeParams
}

type RangeParams struct {
	Gte any `json:"gte,omitempty"`
	Gt  any `json:"gt,omitempty"`
	Lte any `json:"lte,omitempty"`
	Lt  any `json:"lt,omitempty"`
}

func (r *rangeQuery) generate() any {
	rangeParamsMap := map[string]RangeParams{}
	rangeParamsMap[r.fieldName] = r.params
	return struct {
		Range map[string]RangeParams `json:"range"`
	}{
		rangeParamsMap,
	}
}

func Range(field string, params RangeParams) Generatable {
	return &rangeQuery{fieldName: field, params: params}
}

type multiMatchQuery struct {
	params esquery.MultiMatchQuery
}

type MultiMatchParams struct {
	Query  any
	Fields []string
}

func (m *multiMatchQuery) generate() any {
	return m.params.Map()
}

func MultiMatch(params MultiMatchParams) Generatable {
	q := esquery.MultiMatch()
	q.Fields(params.Fields...).Query(params.Query)
	return &multiMatchQuery{params: *q}
}

type aggregationQuery struct {
	params esquery.TermsAggregation
}

type AggregateParams struct {
	Name      string
	FieldName string
	Order     map[string]string
	Size      int
}

func (m *aggregationQuery) generate() any {
	return map[string]map[string]any{
		m.params.Name(): m.params.Map(),
	}
}

func TermsAgg(params AggregateParams) Generatable {
	q := esquery.TermsAgg(params.Name, params.FieldName)
	if params.Size != 0 {
		q.Size(uint64(params.Size))
	}
	if params.Order != nil {
		q.Order(params.Order)
	}

	return &aggregationQuery{params: *q}
}
