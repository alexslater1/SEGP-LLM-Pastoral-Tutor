package storage

type OrderBy string

const (
	OrderByAsc  OrderBy = "asc"
	OrderByDesc OrderBy = "desc"
)

type queryBuilderOrderBy struct {
	column string
	order  OrderBy
}

type QueryBuilder struct {
	limit             *int
	orderBy           *queryBuilderOrderBy
	matchingFields    map[string]string
	greaterThanFields map[string]interface{}
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		matchingFields:    make(map[string]string),
		greaterThanFields: make(map[string]interface{}),
	}
}

func (q *QueryBuilder) OrderBy(column string, order OrderBy) *QueryBuilder {
	q.orderBy = &queryBuilderOrderBy{column: column, order: order}
	return q
}

func (q *QueryBuilder) Limit(limit int) *QueryBuilder {
	q.limit = &limit
	return q
}

func (q *QueryBuilder) Eq(column string, value string) *QueryBuilder {
	q.matchingFields[column] = value
	return q
}

func (q *QueryBuilder) Gt(column string, value interface{}) *QueryBuilder {
	q.greaterThanFields[column] = value
	return q
}
