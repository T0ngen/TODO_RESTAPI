package tasks

import (
	"TodoRESTAPI/internal/http-server/handlers/tasks/mocks"
	"TodoRESTAPI/internal/storage/postgresql"
	"encoding/json"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)



func TestGetAllTasks(t *testing.T) {
    cases := []struct {
        name       string
        want       []postgresql.Task
		basicAuth bool
        respError  int
        mockError  error
    }{
        {
            name: "success",
            want: []postgresql.Task{
                {Id: 1, Text: "good", Importance: 1},
                {Id: 2, Text: "good2", Importance: 2},
            },
			basicAuth: true,
            mockError: nil,
        },
        {
            name: "unsuccess",
            want: nil,
			basicAuth: false,
			respError:  http.StatusInternalServerError,
            mockError: errors.New("authentication failed"),

        },
    }

    for _, tc := range cases {
        tc := tc

        t.Run(tc.name, func(t *testing.T) {
            allTaskskMock := mocks.NewAllTasksk(t)

            handler := All(allTaskskMock)
            

            req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
			
            require.NoError(t, err)
			if tc.basicAuth{
				req.SetBasicAuth("test_user", "password")
				allTaskskMock.On("CheckAllUserTasks", "test_user").Return(tc.want, nil)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				
				require.Equal(t, http.StatusOK, rr.Code)
				
				expected, _ := json.Marshal(tc.want)
				
				require.Equal(t, string(expected), rr.Body.String()) 
	
				allTaskskMock.AssertExpectations(t)
			}
			if !tc.basicAuth{
				allTaskskMock.On("CheckAllUserTasks", "").Return(tc.want, tc.mockError)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				require.Equal(t, tc.respError, rr.Code)
				require.Contains(t, rr.Body.String(), http.StatusText(http.StatusInternalServerError))

				
			}
            
			
           
        })
    }
}