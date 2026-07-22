package data

import "context"

type Inserter[T any] interface{
	Insert(context.Context, T) error
}

type Getter[T any, P any] interface{
	Get(ctx context.Context, id P) (T, error)
}

type Updater[T any] interface{
	Update(context.Context ,T) error
}

type Deleter[T any, P any] interface{
	Delete(ctx context.Context, id P) error
}

// Collection of data models that satisfy the above interface
type Models struct{
	MovieRepository interface{
		Inserter[*Movie]
		Getter[*Movie, int64]
		Updater[*Movie]
		Deleter[*Movie, int64]
	}
}

func New(m MovieModel) Models{
	return Models{
		MovieRepository: m,
	}
}

func NewMock(m MovieModelMock) Models{
	return Models{
		MovieRepository: m,
	}
}