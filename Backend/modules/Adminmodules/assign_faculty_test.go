package Adminmodules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssignFaculty(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Success - Assign faculty to course",
			payload: map[string]interface{}{
				"courseID":   1,
				"facultyID":  3,
				"sectionID":  2,
				"semesterID": 4,
				"deptID":     1,
				"createdBy":  1,
				"updatedBy":  1,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields - Should return 400",
			payload: map[string]interface{}{
				"courseID":  1,
				"facultyID": 3,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid course ID - Should return 400",
			payload: map[string]interface{}{
				"courseID":   999, // Invalid course ID
				"facultyID":  2,
				"sectionID":  3,
				"semesterID": 4,
				"deptID":     5,
				"createdBy":  1,
				"updatedBy":  1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid faculty ID - Should return 400",
			payload: map[string]interface{}{
				"courseID":   1,
				"facultyID":  999, // Invalid faculty ID
				"sectionID":  3,
				"semesterID": 4,
				"deptID":     5,
				"createdBy":  1,
				"updatedBy":  1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Conflict - Faculty already assigned to this course",
			payload: map[string]interface{}{
				"courseID":   1,
				"facultyID":  2,
				"sectionID":  3,
				"semesterID": 4,
				"deptID":     5,
				"createdBy":  1,
				"updatedBy":  1,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert payload to JSON
			body, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "http://localhost:8080/api/admin/assignFaculty", bytes.NewBuffer(body))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Log the response status for debugging
			fmt.Printf("[%s] Status Code: %d\n", tt.name, res.StatusCode)

			// Check if status matches expected value
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
