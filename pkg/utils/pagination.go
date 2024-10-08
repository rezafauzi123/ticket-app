package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Pagination struct {
	Page    int `json:"page" validate:"omitempty,min=1"`
	PerPage int `json:"per_page" validate:"omitempty,min=1,max=100"`
}

func ApplyPagination(query string, pagination Pagination, idx *int) (string, []interface{}) {
	// Cek apakah pagination perlu diterapkan
	if pagination.Page > 0 && pagination.PerPage > 0 {
		limit := pagination.PerPage
		offset := (pagination.Page - 1) * pagination.PerPage
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", *idx, *idx+1)
		return query, []interface{}{limit, offset}
	}

	// Jika pagination tidak diterapkan, kembalikan query dan args tanpa modifikasi
	return query, []interface{}{}
}

func ExtractPagination(r *http.Request, defaultPage, defaultPerPage, maxPerPage int) (Pagination, error) {
	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("per_page")

	var pageInt int
	var perPageInt int
	var err error

	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			return Pagination{}, errors.New("invalid page value")
		}
	} else {
		pageInt = defaultPage // Gunakan nilai default jika tidak diisi
	}

	if perPage != "" {
		perPageInt, err = strconv.Atoi(perPage)
		if err != nil || perPageInt < 1 || perPageInt > maxPerPage {
			return Pagination{}, errors.New("invalid per_page value")
		}
	} else {
		perPageInt = defaultPerPage // Gunakan nilai default jika tidak diisi
	}

	return Pagination{
		Page:    pageInt,
		PerPage: perPageInt,
	}, nil
}
