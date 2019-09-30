package paymentsgateway
import (
	"encoding/json"
	"net/http"
)

func IndexEndpoint(w http.ResponseWriter, req *http.Request) {
	type Response struct {
		Message string
	}
	var resp Response
	resp.Message = "Hello world"
	json.NewEncoder(w).Encode(resp)
}
