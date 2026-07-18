package data

type Mapper[T any] interface {
	Map() T
}

type MovieDto struct {
	Title   string   `json:"title" validate:"required,max=500"`
	Year    int32    `json:"year" validate:"required,yearfrom=1888"`
	Runtime Runtime  `json:"runtime" validate:"required,gt=0"`
	Genres  []string `json:"genres" validate:"required,unique"`
}

func (dto MovieDto) Map() *Movie {
	return &Movie{
		Title:   dto.Title,
		Year:    dto.Year,
		Runtime: dto.Runtime,
		Genres:  dto.Genres,
	}
}

func (movie *Movie) AddId(id int64) *Movie {
	if id > 0 && movie.ID == 0 {
		movie.ID = id
	}
	return movie
}
