package Adminmodules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCourse(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Success - Register a course",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerCourse",
			payload: map[string]interface{}{
				"courseCode":   "CS118",
				"courseName":   "test sub",
				"courseDeptID": 1,
				"courseType":   "L",
				"updatedBy":    1,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Missing Required Fields - Should return 400",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerCourse",
			payload: map[string]interface{}{
				"courseCode":   "",
				"courseName":   "Algorithms",
				"courseDeptID": 1,
				"courseType":   "Core",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid Department ID - Should return 400",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerCourse",
			payload: map[string]interface{}{
				"courseCode":   "CS102",
				"courseName":   "Operating Systems",
				"courseDeptID": 999, // Non-existing deptID
				"courseType":   "Core",
				"updatedBy":    1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid UpdatedBy Admin ID - Should return 400",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerCourse",
			payload: map[string]interface{}{
				"courseCode":   "CS103",
				"courseName":   "Machine Learning",
				"courseDeptID": 1,
				"courseType":   "Elective",
				"updatedBy":    999, // Non-existing admin ID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Method Not Allowed - Should return 405",
			method:         "GET",
			url:            "http://localhost:8080/api/admin/registerCourse",
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
