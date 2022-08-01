package entities

type CompanyCountry struct {
	SequentialIdentifier
	CompanyID       int64               `json:"company_id"`
	CountryID       int64               `json:"country_id"`
	OperationStatus OperationStatusType `json:"operation_status"`
	Timestamps
}

func BuildCompanyCountry(companyID, countryID int64) *CompanyCountry {
	return &CompanyCountry{
		CompanyID:       companyID,
		CountryID:       countryID,
		OperationStatus: OperationStatusTypeActive,
	}
}
