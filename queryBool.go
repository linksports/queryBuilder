package queryBuilder

type boolQuery struct {
	bool boolConditions
}

func (q boolQuery) generate() any {
	return struct {
		Bool boolConditions `json:"bool,omitempty"`
	}{q.bool}
}

type boolConditions struct {
	Must    []any `json:"must,omitempty"`
	MustNot []any `json:"must_not,omitempty"`
	Should  []any `json:"should,omitempty"`
}

func Bool() *boolQuery {
	return &boolQuery{boolConditions{}}
}

func (q *boolQuery) Must(g ...Generatable) *boolQuery {
	must := make([]any, len(g))
	for i, c := range g {
		must[i] = c.generate()
	}
	q.bool.Must = append(q.bool.Must, must...)
	return q
}

func (q *boolQuery) MustNot(g ...Generatable) *boolQuery {
	mustNot := make([]any, len(g))
	for i, c := range g {
		mustNot[i] = c.generate()
	}
	q.bool.MustNot = append(q.bool.MustNot, mustNot...)
	return q
}

func (q *boolQuery) Should(g ...Generatable) *boolQuery {
	should := make([]any, len(g))
	for i, c := range g {
		should[i] = c.generate()
	}
	q.bool.Should = append(q.bool.Should, should...)
	return q
}
