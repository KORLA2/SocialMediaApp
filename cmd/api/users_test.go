package main

import (
	"net/http"
	
	"testing"
)

func TestGetUser(t *testing.T) {
app:=NewTestApplication(t);
	mux:=app.mount();
 t.Run("Test will not allow unauthenticated users", func(t *testing.T) {
	   req,_:=http.NewRequest(http.MethodGet,"/api/v1/user/1",nil)
      rr:= executeRequest(req,mux);

	  if rr.Result().StatusCode!=http.StatusUnauthorized{
		t.Errorf("expected status code %d but got %d",http.StatusUnauthorized,rr.Result().StatusCode)
	  }
	})


}
