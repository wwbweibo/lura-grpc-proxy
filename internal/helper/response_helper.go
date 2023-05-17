package helper

import (
	"encoding/json"
	"net/http"

	"github.com/wwbweibo/lura-grpc-proxy/internal/domain"
)

func WriteResponse(status int, response domain.HttpResponse, writer http.ResponseWriter) {
	bts, _ := json.Marshal(response)
	writer.WriteHeader(status)
	writer.Write(bts)
}
