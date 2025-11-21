package service

import (
	"fmt"
	"kelp-go-assessment/internal/models"
	"log"
	"time"

	"golang.org/x/sync/singleflight"
)

type DataService struct {
	initialDataGroup singleflight.Group
	apiGroup         singleflight.Group
}

func NewDataService() *DataService {
	return &DataService{}
}

// GetInitialData simulates a heavy calculation for initial company data.
// It uses singleflight to ensure it's only calculated once per company concurrently.
func (s *DataService) GetInitialData(companyID string) (*models.InitialData, error) {
	key := fmt.Sprintf("initial_data:%s", companyID)

	v, err, shared := s.initialDataGroup.Do(key, func() (interface{}, error) {
		log.Printf("[InitialData] Starting calculation for Company %s...", companyID)
		// Simulate heavy processing
		time.Sleep(2 * time.Second)
		
		data := &models.InitialData{
			CompanyID:    companyID,
			BaseValue:    1000.0, // Dummy calculation
			Region:       "US-East",
			CalculatedAt: time.Now().Format(time.RFC3339),
		}
		log.Printf("[InitialData] Finished calculation for Company %s", companyID)
		return data, nil
	})

	if err != nil {
		return nil, err
	}

	if shared {
		log.Printf("[InitialData] Shared result for Company %s", companyID)
	}

	return v.(*models.InitialData), nil
}

// GetFinancials calculates financials, reusing InitialData and coalescing concurrent requests.
func (s *DataService) GetFinancials(companyID string) (*models.FinancialsResponse, error) {
	key := fmt.Sprintf("financials:%s", companyID)

	v, err, shared := s.apiGroup.Do(key, func() (interface{}, error) {
		log.Printf("[Financials] Processing request for Company %s...", companyID)
		
		// Step 1: Get Initial Data (this is also coalesced)
		initData, err := s.GetInitialData(companyID)
		if err != nil {
			return nil, err
		}

		// Step 2: Calculate Financials based on Initial Data
		// Simulate some specific processing
		time.Sleep(500 * time.Millisecond)

		resp := &models.FinancialsResponse{
			CompanyID: companyID,
			Revenue:   initData.BaseValue * 10.5,
			Profit:    initData.BaseValue * 2.5,
			Source:    "Calculated",
		}
		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	resp := v.(*models.FinancialsResponse)
	if shared {
		resp.Source = "Coalesced (Shared)"
		log.Printf("[Financials] Returning shared response for Company %s", companyID)
	}
	return resp, nil
}

// GetSales calculates sales data, reusing InitialData and coalescing concurrent requests.
func (s *DataService) GetSales(companyID string) (*models.SalesResponse, error) {
	key := fmt.Sprintf("sales:%s", companyID)

	v, err, shared := s.apiGroup.Do(key, func() (interface{}, error) {
		log.Printf("[Sales] Processing request for Company %s...", companyID)

		initData, err := s.GetInitialData(companyID)
		if err != nil {
			return nil, err
		}

		time.Sleep(500 * time.Millisecond)

		resp := &models.SalesResponse{
			CompanyID:  companyID,
			TotalSales: int(initData.BaseValue * 50),
			Target:     int(initData.BaseValue * 60),
			Source:     "Calculated",
		}
		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	resp := v.(*models.SalesResponse)
	if shared {
		resp.Source = "Coalesced (Shared)"
		log.Printf("[Sales] Returning shared response for Company %s", companyID)
	}
	return resp, nil
}

// GetEmployeeStats calculates employee stats, reusing InitialData and coalescing concurrent requests.
func (s *DataService) GetEmployeeStats(companyID string) (*models.EmployeeStatsResponse, error) {
	key := fmt.Sprintf("employees:%s", companyID)

	v, err, shared := s.apiGroup.Do(key, func() (interface{}, error) {
		log.Printf("[Employees] Processing request for Company %s...", companyID)

		initData, err := s.GetInitialData(companyID)
		if err != nil {
			return nil, err
		}

		time.Sleep(500 * time.Millisecond)

		resp := &models.EmployeeStatsResponse{
			CompanyID: companyID,
			Count:     int(initData.BaseValue / 10),
			AvgTenure: 3.5,
			Source:    "Calculated",
		}
		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	resp := v.(*models.EmployeeStatsResponse)
	if shared {
		resp.Source = "Coalesced (Shared)"
		log.Printf("[Employees] Returning shared response for Company %s", companyID)
	}
	return resp, nil
}
