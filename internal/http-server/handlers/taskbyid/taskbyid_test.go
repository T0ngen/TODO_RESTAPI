package taskbyid

import (
	"TodoRESTAPI/internal/http-server/handlers/taskbyid/mocks"
	"TodoRESTAPI/internal/storage/postgresql"
	"encoding/json"
	"errors"

	// "strconv"

	"fmt"
	
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)






func TestTaskById(t *testing.T){
	cases := []struct{
		name string
		id string
		want *postgresql.TaskById
		respError int
		mockError error
		basicAuth bool
	}{
		{
			name: "success",
			id: "21",
			want: &postgresql.TaskById{Note: "success note",Importance: 2},
			basicAuth: true,
			mockError: nil,
			respError: 200,
		},
		{
			name: "without basicAuth",
			id: "21",
			want: &postgresql.TaskById{Note: "success note",Importance: 2},
			respError: 400,
			basicAuth: false,
			mockError: errors.New("Error: eror with db : handlers.taskbyid.taskbyid.ByID"),
		},
		// {
		// 	name: "without id",
		// 	id: "",
		// 	want: nil,
		// 	respError: 400,
		// 	basicAuth: true,
		// 	mockError: errors.New("Error: empty id : handlers.taskbyid.taskbyid.ByID"),
		// },
	}

	for _, tc := range cases{
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			taskByIdMock := mocks.NewTaskByIdInterface(t)

			handler := ById(taskByIdMock)
			query := "/tasks/?id="+tc.id
			
		
			req:= httptest.NewRequest("GET", query, nil)
			
			
			

			if tc.basicAuth{
				
				req.SetBasicAuth("test_user", "password")
				taskByIdMock.On("CheckTaskById", "test_user", tc.id).Return(tc.want, tc.mockError)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				require.Equal(t, tc.respError, rr.Code)
				
				expected, _ := json.Marshal(tc.want)
				
				
				
				
				require.Equal(t, string(expected), rr.Body.String()) 
	
				taskByIdMock.AssertExpectations(t)
			
			
			}
			if !tc.basicAuth{
			
				taskByIdMock.On("CheckTaskById", "",tc.id).Return(tc.want, tc.mockError)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				fmt.Println(rr.Body.String())
				require.Equal(t, tc.respError, rr.Code)
				
				taskByIdMock.AssertExpectations(t)

			}

		})
	}
	
}


func TestTaskWithoutId(t *testing.T){

	// taskId := ""
	respText := "{\"status\":\"Error\",\"error\":\"empty id\"}\n"
	respError := 400

	taskByIdMock := mocks.NewTaskByIdInterface(t)

	handler := ById(taskByIdMock)
	query := "/tasks/?id="
	req:= httptest.NewRequest("GET", query, nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, respError, rr.Code)
	require.Equal(t, respText, rr.Body.String())
	

}

func TestTaskResultNil(t *testing.T){
	username := "qwerty"
	respText :="{\"status\":\"Error\",\"error\":\"invalid request\"}\n"
	respError := 400
	taskId := "1"
	taskByIdMock := mocks.NewTaskByIdInterface(t)
	taskByIdMock.On("CheckTaskById", username, taskId).Return(nil, nil)
	handler := ById(taskByIdMock)

	query := "/tasks/?id=" + taskId
	req:= httptest.NewRequest("GET", query, nil)
	req.SetBasicAuth(username, "password")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	

	require.Equal(t, respText, rr.Body.String())
	require.Equal(t, respError, rr.Code)

}