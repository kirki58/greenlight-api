package data

type Inserter[T any] interface{
	Insert(T) error
}

type Getter[T any, P any] interface{
	Get(id P) (T, error)
}

type Updater[T any] interface{
	Update(T) error
}

type Deleter[T any, P any] interface{
	Delete(id P) error
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