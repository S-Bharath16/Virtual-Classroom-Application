package Adminmodules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCourse(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Success - Update a course",
			method: "PUT",
			url:    "http://localhost:8080/api/admin/updateCourse",
			payload: map[string]interface{}{
				"courseID":     1,
				"courseCode":   "CS202",
				"courseName":   "Advanced Algorithms",
				"courseDeptID": 1,
				"courseType":   "L",
				"updatedBy":    1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Missing Required Fields - Should return 400",
			method: "PUT",
			url:    "http://localhost:8080/api/admin/updateCourse",
			payload: map[string]interface{}{
				"courseID":     1,
				"courseCode":   "",
				"courseName":   "Algorithms",
				"courseDeptID": 1,
				"courseType":   "Core",
				"updatedBy":    1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid Course ID - Should return 400",
			method: "PUT",
			url:    "http://localhost:8080/api/admin/updateCourse",
			payload: map[string]interface{}{
				"courseID":     999, // Non-existing course ID
				"courseCode":   "CS202",
				"courseName":   "Advanced Algorithms",
				"courseDeptID": 1,
				"courseType":   "L",
				"updatedBy":    1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid Department ID - Should return 400",
			method: "PUT",
			url:    "http://localhost:8080/api/admin/updateCourse",
			payload: map[string]interface{}{
				"courseID":     1,
				"courseCode":   "CS202",
				"courseName":   "Advanced Algorithms",
				"courseDeptID": 999, // Invalid department ID
				"courseType":   "L",
				"updatedBy":    1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid UpdatedBy Admin ID - Should return 400",
			method: "PUT",
			url:    "http://localhost:8080/api/admin/updateCourse",
			payload: map[string]interface{}{
				"courseID":     1,
				"courseCode":   "CS202",
				"courseName":   "Advanced Algorithms",
				"courseDeptID": 1,
				"courseType":   "L",
				"updatedBy":    999, // Non-existing admin ID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Method Not Allowed - Should return 405",
			method:         "GET",
			url:            "http://localhost:8080/api/admin/updateCourse",
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
