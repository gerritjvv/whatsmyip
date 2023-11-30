package internal

import (
	"net/http"
	"strings"
)

type MyIpResponse struct {
	Headers       map[string]string `json:"headers"`
	RemoteAddress string            `json:"remote_address"`
}

func ReturnMyIpHandler(request *http.Request) MyIpResponse {

	headers := make(map[string]string)
	i := 0
	max_headers := 20
	for h, v := range request.Header {
		if i >= max_headers {
			break
		}
		if len(v) > 0 {
			if len(v) > 5 {
				v = v[:5]
			}
			s := strings.Join(v, ",")
			headers[h] = s
		}
	}

	remoteAdd := request.RemoteAddr

	return MyIpResponse{
		Headers:       headers,
		RemoteAddress: remoteAdd,
	}
}
