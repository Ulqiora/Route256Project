package shared

import (
	"encoding/json"
	"net/http"
)

func WriteToWriter(response http.ResponseWriter, a any) {
	bytesAll, err := json.Marshal(a)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = response.Write(bytesAll)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}
