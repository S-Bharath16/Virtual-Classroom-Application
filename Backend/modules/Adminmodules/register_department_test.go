package Adminmodules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterDepartment(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Success - Register a new department",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerDepartment",
			payload: map[string]interface{}{
				"deptName": "ECE",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Missing Required Field - Should return 400",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerDepartment",
			payload: map[string]interface{}{
				"deptName": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Method Not Allowed - Should return 405",
			method:         "GET",
			url:            "http://localhost:8080/api/admin/registerDepartment",
			payload:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if tt.payload != nil {
				reqBody, err = json.Marshal(tt.payload)
				assert.NoError(t, err)
			}

			req, err := http.NewRequest(tt.method, tt.url, bytes.NewBuffer(reqBody))
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
