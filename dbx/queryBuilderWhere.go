package dbx

func (q *QrBuilder[T]) Where(where string, args ...interface{}) *QrBuilder[T] {
	q.where = where
	q.args = args
	return q
}
func (q *QrBuilder[T]) Select(selector string) *QrBuilder[T] {
	q.selector = selector
	return q
}
