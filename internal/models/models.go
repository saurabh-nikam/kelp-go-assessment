package models

// InitialData represents the shared data calculated once per company
type InitialData struct {
	CompanyID   string
	BaseValue   float64
	Region      string
	CalculatedAt string
}

// FinancialsResponse represents the response for the financials API
type FinancialsResponse struct {
	CompanyID string  `json:"company_id"`
	Revenue   float64 `json:"revenue"`
	Profit    float64 `json:"profit"`
	Source    string  `json:"source"` // To indicate if it was coalesced or fresh
}

// SalesResponse represents the response for the sales API
type SalesResponse struct {
	CompanyID string  `json:"company_id"`
	TotalSales int    `json:"total_sales"`
	Target    int     `json:"target"`
	Source    string  `json:"source"`
}

// EmployeeStatsResponse represents the response for the employee stats API
type EmployeeStatsResponse struct {
	CompanyID string `json:"company_id"`
	Count     int    `json:"count"`
	AvgTenure float64 `json:"avg_tenure"`
	Source    string `json:"source"`
}
