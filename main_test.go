package queryBuilder_test

import (
	"testing"

	queryBuilder "github.com/linksports/esQueryBuilder"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	t.Run("match_all", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.MatchAll(),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"match_all":{}
			}
		}`), query)
	})

	t.Run("match", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Match("target", "v"),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"match":{"target":"v"}
			}
		}`), query)
	})

	t.Run("match_phrase", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.MatchPhrase("target", []string{"red", "blue", "green"}),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"match_phrase":{"target":"red blue green"}
			}
		}`), query)
	})

	t.Run("term", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Term("target.keyword", "v"),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"term":{"target.keyword":"v"}
			}
		}`), query)
	})

	t.Run("terms(int)", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Terms("target.keyword", []int{1, 2, 3}),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"terms":{"target.keyword":[1,2,3]}
			}
		}`), query)
	})

	t.Run("terms(string)", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Terms("target.keyword", []string{"1", "2", "3"}),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"terms":{"target.keyword":["1","2","3"]}
			}
		}`), query)
	})

	t.Run("prefix", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Prefix("target", "v"),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"prefix":{"target":"v"}
			}
		}`), query)
	})

	t.Run("query+sort", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Term("target.keyword", "v"),
		).Sort(
			queryBuilder.Sort{"sort1", "asc"},
			queryBuilder.Sort{"sort2", "desc"},
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"term":{"target.keyword":"v"}
			},
			"sort":[
				{"sort1":{"order":"asc"}},
				{"sort2":{"order":"desc"}}
			]
		}`), query)
	})

	t.Run("query+_source", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Term("target.keyword", "v"),
		).Source([]string{"a", "b"}).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"term":{"target.keyword":"v"}
			},
			"_source":["a","b"]
		}`), query)
	})

	t.Run("query+sort+size+from+_source", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Term("target.keyword", "v"),
		).Sort(
			queryBuilder.Sort{"sort1", "asc"},
		).Source([]string{"taro", "hanako"}).Size(10).From(5).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"size":10,
			"from":5,
			"query":{
				"term":{"target.keyword":"v"}
			},
			"sort":[
				{"sort1":{"order":"asc"}}
			],
			"_source":["taro","hanako"]
		}`), query)
	})

	t.Run("aggs_term", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Aggs(
			queryBuilder.TermsAgg(queryBuilder.AggregateParams{
				Name:      "sportID_term",
				FieldName: "sport_id",
				Order: map[string]string{
					"_count": "desc",
				},
				Size: 10,
			}),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"aggs":{
				"sportID_term":{
					"terms":{
						"field":"sport_id",
						"order":{"_count":"desc"},
						"size":10
						}
				}
			}
		}`), query)
	})
	t.Run("multi_aggs_term", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Aggs(
			queryBuilder.TermsAgg(queryBuilder.AggregateParams{
				Name:      "sportID_term",
				FieldName: "sport_id",
				Order: map[string]string{
					"_count": "desc",
				},
				Size: 10,
			}),
			queryBuilder.TermsAgg(queryBuilder.AggregateParams{
				Name:      "sportID_term2",
				FieldName: "sport_id2",
				Order: map[string]string{
					"_count": "asc",
				},
				Size: 12,
			}),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"aggs":{
				"sportID_term":{
					"terms":{
						"field":"sport_id",
						"order":{"_count":"desc"},
						"size":10
						}
				},
				"sportID_term2":{
					"terms":{
						"field":"sport_id2",
						"order":{"_count":"asc"},
						"size":12
					}
				}
			}
		}`), query)
	})

	t.Run("query+aggs_term+size+from", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.FunctionScore(queryBuilder.Bool().Must(
				queryBuilder.MultiMatch(queryBuilder.MultiMatchParams{
					Query:  "teamName1 teamName2 tokyo",
					Fields: []string{"name^3", "description", "city"},
				}),
			).MustNot(
				queryBuilder.Exists("deleted_at"),
			),
				[]queryBuilder.Function{
					{
						Filter: queryBuilder.Exists("logo"),
						Weight: 3,
					},
					{
						Filter: queryBuilder.Exists("photo"),
						Weight: 1.5,
					},
				},
			),
		).Aggs(
			queryBuilder.TermsAgg(queryBuilder.AggregateParams{
				Name:      "sportID_term",
				FieldName: "sport_id",
				Order: map[string]string{
					"_count": "desc",
				},
				Size: 10,
			}),
		).Size(20).From(5).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{"size":20,
			"from":5,
			"query":{
				"function_score":{
					"query":{
						"bool":{
							"must":[{
								"multi_match":{
									"fields":["name^3","description","city"],
									"query":"teamName1 teamName2 tokyo"
								}	
							}],
							"must_not":[{
								"exists":{
									"field":"deleted_at"
								}
							}]
						}
					},
					"functions":[
						{
							"filter":{
								"exists":{
									"field":"logo"
								}
							},
							"weight":3
						},
						{
							"filter":{
								"exists":{
									"field":"photo"
								}
							},
							"weight":1.5
						}
					]
				}
			},
			"aggs":{
				"sportID_term":{
					"terms":{
						"field":"sport_id",
						"order":{"_count":"desc"},
						"size":10
					}
				}
			}
		}`), query)
	})

	t.Run("search_after", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.SearchAfter("0", "1").Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"search_after":["0","1"]
		}`), query)
	})

	t.Run("exists", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Exists("target"),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"exists":{"field":"target"}
			}
		}`), query)
	})

	t.Run("range", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.Range("target", queryBuilder.RangeParams{
				Gte: 10,
				Lt:  "hanako",
			}),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"range":{
					"target":{
						"gte":10,
						"lt":"hanako"
					}
				}
			}
		}`), query)
	})

	t.Run("multi_match", func(t *testing.T) {
		builder := queryBuilder.New()
		query, err := builder.Query(
			queryBuilder.MultiMatch(queryBuilder.MultiMatchParams{
				Query:  "Elastic Search",
				Fields: []string{"clusterName", "searchEngine^2"},
			}),
		).Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"multi_match":{
					"fields":[
						"clusterName",
						"searchEngine^2"
					],
					"query":"Elastic Search"
				}
			}
		}`), query)
	})

	t.Run("function_score", func(t *testing.T) {
		builder := queryBuilder.New()
		builder.Query(
			queryBuilder.FunctionScore(
				queryBuilder.Term("target.keyword", "v"),
				[]queryBuilder.Function{
					{
						Filter: queryBuilder.Term("target1", "function1"),
						Weight: 2,
					},
					{
						Filter: queryBuilder.Term("target2", "function2"),
						Weight: 5,
					},
				},
			),
		)

		query, err := builder.Build()

		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"function_score":{
					"query":{
						"term":{"target.keyword":"v"}
					},
					"functions":[
						{
							"filter":{"term":{"target1":"function1"}},
							"weight":2
						},
						{
							"filter":{"term":{"target2":"function2"}},
							"weight":5
						}
					]
				}
			}
		}`), query)
	})
}
