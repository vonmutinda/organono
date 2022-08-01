package forms

import "gopkg.in/guregu/null.v3"

type CreateCompanyForm struct {
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	CountryID int64  `json:"country_id" binding:"required"`
	Website   string `json:"website" bidning:"required"`
	Phone     string `json:"phone" binding:"required"`
}

type UpdateCompanyForm struct {
	Name    null.String `json:"name" binding:"required"`
	Code    null.String `json:"code" binding:"required"`
	Website null.String `json:"website" bidning:"required"`
	Phone   null.String `json:"phone" binding:"required"`
}
