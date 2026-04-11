package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

type PageParams struct {
	Page    int
	PerPage int
}

func ParsePageParams(c *gin.Context) PageParams {
	page := DefaultPage
	perPage := DefaultPerPage

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if pp, err := strconv.Atoi(c.Query("per_page")); err == nil && pp > 0 {
		perPage = pp
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return PageParams{Page: page, PerPage: perPage}
}

func Paginate(params PageParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (params.Page - 1) * params.PerPage
		return db.Offset(offset).Limit(params.PerPage)
	}
}

func NewMeta(params PageParams, total int64) *Meta {
	totalPages := int(math.Ceil(float64(total) / float64(params.PerPage)))
	return &Meta{
		Page:       params.Page,
		PerPage:    params.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
