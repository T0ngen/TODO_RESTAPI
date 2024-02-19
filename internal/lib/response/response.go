package response




type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type ResponseOK struct {
	Status string `json:"status"`
	OK  string `json:"response,omitempty"`
}

func Error(msg string) Response {
	return Response{
		Status: "Error",
		Error:  msg,
	}
}

func RespOK(msg string) ResponseOK {
	return ResponseOK{
		Status: "success",
		OK:  msg,
	}
}