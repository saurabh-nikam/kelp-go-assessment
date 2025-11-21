package test

import (
	"encoding/json"
	"kelp-go-assessment/internal/api/routes"
	"kelp-go-assessment/internal/models"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestRequestCoalescing(t *testing.T) {
	router := routes.SetupRouter()

	// Test Case 1: Concurrent requests to the same endpoint (Financials)
	t.Run("Concurrent Financials", func(t *testing.T) {
		var wg sync.WaitGroup
		numRequests := 5
		responses := make([]models.FinancialsResponse, numRequests)

		start := time.Now()

		for i := 0; i < numRequests; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/api/financials?companyId=100", nil)
				router.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("Request %d failed with code %d", index, w.Code)
					return
				}

				var resp models.FinancialsResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				responses[index] = resp
			}(i)
		}

		wg.Wait()
		duration := time.Since(start)

		// Verification
		// 1. Duration should be around 2.5s (2s initial + 0.5s calculation), not 2.5s * 5
		// But since we are using httptest, it might be slightly different, but definitely not sequential.
		t.Logf("Total duration for %d requests: %v", numRequests, duration)

		calculatedCount := 0
		coalescedCount := 0
		for i, resp := range responses {
			t.Logf("Req %d Source: %s", i, resp.Source)
			switch resp.Source {
			case "Calculated":
				calculatedCount++
			case "Coalesced (Shared)":
				coalescedCount++
			default:
				t.Errorf("Unexpected source: %s", resp.Source)
			}
		}

		// Since singleflight returns 'shared=true' for all callers if multiple are in flight,
		// and we update the Source based on that, it's possible/likely that ALL requests
		// get marked as "Coalesced (Shared)".
		// The important thing is that we didn't get 5 separate calculations (which would take 5 * 2.5s).
		// And we want to see at least some coalescing.

		if coalescedCount == 0 {
			t.Errorf("Expected some coalescing, got 0. All were %d Calculated", calculatedCount)
		}
		if calculatedCount+coalescedCount != numRequests {
			t.Errorf("Total responses mismatch")
		}
		t.Logf("Calculated: %d, Coalesced: %d", calculatedCount, coalescedCount)
	})

	// Test Case 2: Concurrent requests to DIFFERENT endpoints for SAME company
	// They should share the Initial Data calculation.
	t.Run("Shared Initial Data", func(t *testing.T) {
		var wg sync.WaitGroup

		// We will hit Financials and Sales at the same time.
		// Both take 2s (Initial) + 0.5s (Specific).
		// Total time should be around 2.5s, and Initial Data should be computed once.
		// We can't easily check internal logs here, but we can check timing.

		start := time.Now()
		wg.Add(2)

		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/financials?companyId=200", nil)
			router.ServeHTTP(w, req)
		}()

		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/sales?companyId=200", nil)
			router.ServeHTTP(w, req)
		}()

		wg.Wait()
		duration := time.Since(start)
		t.Logf("Total duration for mixed requests: %v", duration)

		// If they were sequential or didn't share initial data, it might take longer or at least we know they finished.
		// The real proof is in the logs, but for this test we ensure it completes successfully and reasonably fast.
		if duration > 4500*time.Millisecond {
			t.Errorf("Took too long (%v), suggesting no coalescing of initial data", duration)
		}
	})
}
