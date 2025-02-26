package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	query := NewQueryBuilder()
	query.Eq("name", "test")

	assert.Equal(t, query.matchingFields, map[string]string{"name": "test"})
}

func TestQueryBuilderOrderBy(t *testing.T) {
	query := NewQueryBuilder()
	query.OrderBy("name", OrderByAsc)
	assert.Equal(t, query.orderBy, &queryBuilderOrderBy{column: "name", order: OrderByAsc})
}

func TestQueryBuilderLimit(t *testing.T) {
	query := NewQueryBuilder()
	query.Limit(10)

	n := 10
	assert.Equal(t, query.limit, &n)
}

func TestComplexQuery(t *testing.T) {
	query := NewQueryBuilder()
	query.Eq("name", "test").Eq("age", "20").OrderBy("name", OrderByAsc).Limit(10)

	assert.Equal(t, query.matchingFields, map[string]string{"name": "test", "age": "20"})
	assert.Equal(t, query.orderBy, &queryBuilderOrderBy{column: "name", order: OrderByAsc})
	assert.Equal(t, *query.limit, 10)
}
