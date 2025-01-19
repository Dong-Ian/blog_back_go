package response

import (
	"encoding/json"
	"net/http"
)

func Response(res http.ResponseWriter, message interface{}) {
	responseObject, _ := json.Marshal(message)

	res.WriteHeader(200)
	res.Write(responseObject)
}
