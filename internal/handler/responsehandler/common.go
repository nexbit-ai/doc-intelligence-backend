package responsehandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"syscall"
)

func GenericRes(w http.ResponseWriter, r *http.Request) {
	resData := r.Context().Value("resData")

	var payload = map[string]interface{}{
		"status": true,
		"error":  "",
		"data":   resData,
	}
	respondwithJSON(w, r, http.StatusOK, payload)
}

func respondwithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if errors.Is(err, syscall.EPIPE) {

		return
	}
	if err != nil {
		panic(err)
	}
}
