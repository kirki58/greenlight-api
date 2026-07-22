package data

type Mapper[T any] interface {
	MapTo(T) T // Map to an existing one (requires pass by reference)
}

type MovieDto struct {
	Title   string   `json:"title" validate:"required,min=0,max=500"`
	Year    int32    `json:"year" validate:"required,yearfrom=1888"`
	Runtime Runtime  `json:"runtime" validate:"required,gt=0"`
	Genres  []string `json:"genres" validate:"required,omitempty,unique,dive,min=0,max=100"`
}

func (dto MovieDto) MapTo(mov *Movie) *Movie {
	return &Movie{
		Title:   dto.Title,
		Year:    dto.Year,
		Runtime: dto.Runtime,
		Genres:  dto.Genres,
	}
}

type PartialMovieDto struct{
	Title   *string   `json:"title" validate:"omitempty,min=0,max=500"`
	Year    *int32    `json:"year" validate:"omitempty,yearfrom=1888"`
	Runtime *Runtime  `json:"runtime" validate:"omitempty,gt=0"`
	Genres  []string `json:"genres" validate:"omitempty,unique,dive,min=0,max=100"`
}

func (dto PartialMovieDto) MapTo(mov *Movie) *Movie {
	if dto.Title != nil{
		mov.Title = *dto.Title
	}
	if dto.Year != nil{
		mov.Year = *dto.Year
	}
	if dto.Runtime != nil{
		mov.Runtime = *dto.Runtime
	}
	if dto.Genres != nil{
		mov.Genres = dto.Genres
	}
	return mov
}
