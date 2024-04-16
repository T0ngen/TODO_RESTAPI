package addtask

import (
	"TodoRESTAPI/internal/http-server/handlers/addtask/mocks"
	"bytes"
	"encoding/json"
	"errors"

	"net/http/httptest"

	"testing"

	"github.com/stretchr/testify/require"
)

// {"note":"Make a film", "importance":1}

func TestAddTask(t *testing.T){
	cases := []struct{
		name string
		username string
		respError int
		success bool
		mockError error
		jsonData string
		respText string
	}{
		{name: "success",
		username: "test_user",
		respError: 200,
		success: true,
		mockError: nil,
		respText:"{\"status\":\"success\",\"response\":\"added to notes\"}\n",
		jsonData: `{"note": "Make a film", "importance": 1}`,	},


		{name: "request body empty",
		username: "test_user",
		respError: 400,
		success: false,
		mockError: nil,
		respText:"{\"status\":\"Error\",\"error\":\"empty request\"}\n",
		jsonData: "",	},

		{name: "error with db",
		username: "test_user",
		respError: 500,
		success: false,
		mockError: errors.New("error on db"),
		respText:"Internal Server Error\n",
		jsonData: `{"note": "Make a film", "importance": 1}`,	},
		{name: "not success after db",
		username: "test_user",
		respError: 500,
		success: false,
		mockError: nil,
		respText:"Internal Server Error\n",
		jsonData: `{"note": "Make a film", "importance": 1}`,	},
	}

	for _, tc := range cases{
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			addTaskMock := mocks.NewAddTaskInterface(t)

			handler := New(addTaskMock)
			
			
			req:= httptest.NewRequest("GET", "/tasks", bytes.NewBuffer([]byte(tc.jsonData)))
			req.SetBasicAuth(tc.username, "password")
			var request Request
			err := json.Unmarshal([]byte(tc.jsonData), &request)
			if err != nil{
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
					
				require.Equal(t, tc.respError, rr.Code)
				require.Equal(t, tc.respText, rr.Body.String())
				return
			}

			addTaskMock.On("AddNewTask", tc.username, request.Note, request.Importance).Return(tc.success, tc.mockError)
			
			

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
				
			require.Equal(t, tc.respError, rr.Code)
			require.Equal(t, tc.respText, rr.Body.String())
			
			addTaskMock.AssertExpectations(t)
	
		})
}
}





func TestAddTaskEmptyBody(t *testing.T){
	cases := []struct{
		name string
		username string
		respError int
		success bool
		mockError error
		jsonData string
		respText string
	}{
		
		{name: "request body empty",
		username: "test_user",
		respError: 400,
		success: false,
		mockError: nil,
		respText:"{\"status\":\"Error\",\"error\":\"empty request\"}\n",
		jsonData: "",	},
	}

	for _, tc := range cases{
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			addTaskMock := mocks.NewAddTaskInterface(t)

			handler := New(addTaskMock)
			
			
			req:= httptest.NewRequest("GET", "/tasks", bytes.NewBuffer([]byte(tc.jsonData)))
			req.SetBasicAuth(tc.username, "password")
			var request Request
			err :=json.Unmarshal([]byte(tc.jsonData), &request)
			// fmt.Println("req",request.Note)
			// fmt.Println("req",request.Importance)
			// require.NoError(t, err)
			if err != nil{
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
					
				require.Equal(t, tc.respError, rr.Code)
				return
			}
			addTaskMock.On("AddNewTask", tc.username, request.Note, request.Importance).Return(tc.success, tc.mockError)
			
			

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
				
			require.Equal(t, tc.respError, rr.Code)
			require.Equal(t, tc.respText, rr.Body.String())
			addTaskMock.AssertExpectations(t)
	
		})
}
}