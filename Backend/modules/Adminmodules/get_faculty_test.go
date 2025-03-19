package Adminmodules

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllFaculty(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "Success - Get all faculty data",
			method:         "GET",
			url:            "http://localhost:8080/api/admin/getAllFaculty",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Endpoint - Should return 404",
			method:         "GET",
			url:            "http://localhost:8080/api/admin/invalidEndpoint",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Method Not Allowed - Should return 405",
			method:         "POST",
			url:            "http://localhost:8080/api/admin/getAllFaculty",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			res, err := client.Do(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Log the response status for debugging
			fmt.Printf("[%s] Status Code: %d\n", tt.name, res.StatusCode)

			// Assert the status code matches the expected status
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
