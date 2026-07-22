package main

type AllMoviesQuery struct {
	Title    string   `schema:"title" validate:"max=500"`
	Genres   []string `schema:"genres" validate:"omitempty,unique,dive,min=0,max=100"`
	Page     int      `schema:"page,default:1" validate:"gt=0,lte=100"`
	PageSize int      `schema:"pageSize,default:10" validate:"gt=0,lte=100"`
	Sort     string   `schema:"sort,default:title" validate:"oneof=title -title runtime -runtime year -year"`
}
