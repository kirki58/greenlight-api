package dto

import "github.com/kirki58/greenlight/m/internal/data"

type Mapper[T any] interface{
	Map() T
}

type MovieDto struct{
	Title   string       `json:"title" validate:"required,max=500"`
	Year    int32        `json:"year" validate:"required,yearfrom=1888"`
	Runtime data.Runtime `json:"runtime" validate:"required,gt=0"`
	Genres  []string     `json:"genres" validate:"required,unique"`
}

func (dto MovieDto) Map() data.Movie{
	return data.Movie{
		Title: dto.Title,
		Year: dto.Year,
		Runtime: dto.Runtime,
		Genres: dto.Genres,
	}
}