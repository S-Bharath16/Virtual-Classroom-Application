package Adminmodules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterStudent(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Success - Register a student",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerStudent",
			payload: map[string]interface{}{
				"rollNumber":  "21BCS123",
				"emailID":     "student@example.com",
				"studentName": "John Doe",
				"startYear":   2021,
				"endYear":     2025,
				"deptID":      1,
				"sectionID":   1,
				"semesterID":  1,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Missing Required Fields - Should return 400",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerStudent",
			payload: map[string]interface{}{
				"rollNumber":  "",
				"emailID":     "student@example.com",
				"studentName": "John Doe",
				"startYear":   2021,
				"endYear":     2025,
				"deptID":      1,
				"sectionID":   1,
				"semesterID":  1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid Department ID - Should return 400",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerStudent",
			payload: map[string]interface{}{
				"rollNumber":  "21BCS124",
				"emailID":     "student2@example.com",
				"studentName": "Jane Doe",
				"startYear":   2021,
				"endYear":     2025,
				"deptID":      999, // Invalid department ID
				"sectionID":   1,
				"semesterID":  1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid Section ID - Should return 400",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerStudent",
			payload: map[string]interface{}{
				"rollNumber":  "21BCS125",
				"emailID":     "student3@example.com",
				"studentName": "Sam Smith",
				"startYear":   2021,
				"endYear":     2025,
				"deptID":      1,
				"sectionID":   999, // Invalid section ID
				"semesterID":  1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid Semester ID - Should return 400",
			method: "POST",
			url:    "http://localhost:8080/api/admin/registerStudent",
			payload: map[string]interface{}{
				"rollNumber":  "21BCS126",
				"emailID":     "student4@example.com",
				"studentName": "Emma Wilson",
				"startYear":   2021,
				"endYear":     2025,
				"deptID":      1,
				"sectionID":   1,
				"semesterID":  999, // Invalid semester ID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Method Not Allowed - Should return 405",
			method:         "GET",
			url:            "http://localhost:8080/api/admin/registerStudent",
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
