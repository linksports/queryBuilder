package queryBuilder_test

import (
	"testing"

	"github.com/linksports/queryBuilder"
	"github.com/stretchr/testify/assert"
)

func TestQueryBool(t *testing.T) {
	t.Run("check query", func(t *testing.T) {
		t.Run("bool > must > term", func(t *testing.T) {
			builder := queryBuilder.New()
			query, err := builder.Query(
				queryBuilder.Bool().Must(
					queryBuilder.Term("target.keyword", "v"),
				),
			).Build(queryBuilder.ES)

			assert.NoError(t, err)
			assert.Equal(t, queryBuilder.Trim(`{
				"query":{
					"bool":{
						"must":[
							{"term":{"target.keyword":"v"}}
						]
					}
				}
			}`), query)
		})

		t.Run("bool > must > term(x2)", func(t *testing.T) {
			builder := queryBuilder.New()
			query, err := builder.Query(
				queryBuilder.Bool().Must(
					queryBuilder.Term("target1.keyword", "v1"),
					queryBuilder.Term("target2.keyword", "v2"),
				),
			).Build(queryBuilder.ES)

			assert.NoError(t, err)
			assert.Equal(t, queryBuilder.Trim(`{
				"query":{
					"bool":{
						"must":[
							{"term":{"target1.keyword":"v1"}},
							{"term":{"target2.keyword":"v2"}}
						]
					}
				}
			}`), query)
		})

		t.Run("bool > must > terms", func(t *testing.T) {
			builder := queryBuilder.New()
			query, err := builder.Query(
				queryBuilder.Bool().Must(
					queryBuilder.Terms("target.keyword", []int{1, 2, 3}),
				),
			).Build(queryBuilder.ES)

			assert.NoError(t, err)
			assert.Equal(t, queryBuilder.Trim(`{
				"query":{
					"bool":{
						"must":[
							{"terms":{"target.keyword":[1,2,3]}}
						]
					}
				}
			}`), query)
		})

		t.Run("bool > must+must_not", func(t *testing.T) {
			builder := queryBuilder.New()
			query, err := builder.Query(
				queryBuilder.Bool().Must(
					queryBuilder.Terms("target.keyword", []int{1, 2, 3}),
				).MustNot(
					queryBuilder.Terms("not.keyword", []string{"1", "2", "3"}),
				),
			).Build(queryBuilder.ES)

			assert.NoError(t, err)
			assert.Equal(t, queryBuilder.Trim(`{
				"query":{
					"bool":{
						"must":[
							{"terms":{"target.keyword":[1,2,3]}}
						],
						"must_not":[
							{"terms":{"not.keyword":["1","2","3"]}}
						]
					}
				}
			}`), query)
		})

		t.Run("bool > should > bool > must", func(t *testing.T) {
			builder := queryBuilder.New()
			query, err := builder.Query(
				queryBuilder.Bool().Should(
					queryBuilder.Bool().Must(
						queryBuilder.Terms("target1.keyword", []int{1, 2, 3}),
						queryBuilder.Terms("target2.keyword", []int{1, 2, 3}),
					),
					queryBuilder.Bool().Must(
						queryBuilder.Terms("target3.keyword", []int{1, 2, 3}),
						queryBuilder.Terms("target4.keyword", []int{1, 2, 3}),
					),
				),
			).Build(queryBuilder.ES)

			assert.NoError(t, err)
			assert.Equal(t, queryBuilder.Trim(`{
				"query":{
					"bool":{
						"should":[
							{
								"bool":{
									"must":[
										{"terms":{"target1.keyword":[1,2,3]}},
										{"terms":{"target2.keyword":[1,2,3]}}
									]
								}
							},
							{
								"bool":{
									"must":[
										{"terms":{"target3.keyword":[1,2,3]}},
										{"terms":{"target4.keyword":[1,2,3]}}
									]
								}
							}
						]
					}
				}
			}`), query)
		})
	})

	t.Run("building query", func(t *testing.T) {
		boolQuery := queryBuilder.Bool()
		boolQuery.Must(queryBuilder.Term("term", "value"))
		boolQuery.Must(queryBuilder.Terms("terms", []int{1, 2, 3}))
		boolQuery.MustNot(queryBuilder.Term("notTarget", "exclusion"))

		query, err := queryBuilder.New().Query(boolQuery).Build(queryBuilder.ES)
		assert.NoError(t, err)
		assert.Equal(t, queryBuilder.Trim(`{
			"query":{
				"bool":{
					"must":[
						{"term":{"term":"value"}},
						{"terms":{"terms":[1,2,3]}}
					],
					"must_not":[
						{"term":{"notTarget":"exclusion"}}
					]
				}
			}
		}`), query)
	})
}
