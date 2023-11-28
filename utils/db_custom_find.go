package utils

import "gorm.io/gorm"

type ReadOption struct {
	Limit    *int
	Offset   *int
	Order    *string
	Where    map[string]interface{}
	Many     bool
	Preloads []string
	Result   interface{}
}

func (r *ReadOption) Read(db *gorm.DB) *gorm.DB {
	stmt := db
	if r.Where != nil {
		stmt = stmt.Where(r.Where)
	}
	if r.Limit != nil {
		stmt = stmt.Limit(*r.Limit)
	}
	if r.Offset != nil {
		stmt = stmt.Offset(*r.Offset)
	}
	if r.Order != nil {
		stmt = stmt.Order(*r.Order)
	}
	for _, preload := range r.Preloads {
		stmt = stmt.Preload(preload)
	}
	if r.Many {
		return stmt.Find(r.Result)
	} else {
		return stmt.First(r.Result)
	}
}
