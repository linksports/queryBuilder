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
