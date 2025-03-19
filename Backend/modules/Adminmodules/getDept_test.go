package Adminmodules

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetDept(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	} {
		{
			name:           "Success - Get all departments",
			method:         "GET",
			url:            "http://localhost:8080/api/admin/getAllDepartments",
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
			url:            "http://localhost:8080/api/admin/getAllDepartments",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			assert.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}