package schema

type Query[T any] interface {
	// SELECT * FROM @@table WHERE id=@id AND deleted_at=0
	FindOne(id int64) (T, error)

	// SELECT COUNT(1) FROM @@table WHERE id=@id AND deleted_at=0
	Exist(id int64) (int64, error)
}
