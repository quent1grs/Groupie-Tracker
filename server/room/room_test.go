package room

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRoom(t *testing.T) {
	req, err := http.NewRequest("GET", "/room/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRoom)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}

func TestHandleCreateRoom(t *testing.T) {
	roomRequest := `{
        "gameName": "Test Game",
        "gameType": "Test Type",
        "rounds": "5",
        "roundDuration": "10",
        "maxPlayers": "4",
        "scattegories": ["cat1", "cat2", "cat3"]
    }`
	req, err := http.NewRequest("POST", "/room", strings.NewReader(roomRequest))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleCreateRoom)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}
