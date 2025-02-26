package storage

type OrderBy string

const (
	OrderByAsc  OrderBy = "ASC"
	OrderByDesc OrderBy = "DESC"
)

type queryBuilderOrderBy struct {
	column string
	order  OrderBy
}

type QueryBuilder struct {
	limit          *int
	orderBy        *queryBuilderOrderBy
	matchingFields map[string]string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		matchingFields: make(map[string]string),
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
