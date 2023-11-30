package main

import (
	"encoding/json"
	"fmt"
	"github.com/gerritjvv/whatsmyip/internal"
	"github.com/gorilla/mux"
	"net/http"
)

type RespErr struct {
	Message string `json:"message"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/myip", func(writer http.ResponseWriter, request *http.Request) {
		resp := internal.ReturnMyIpHandler(request)
		respBytes, err := json.Marshal(resp)
		if err != nil {
			writer.WriteHeader(500)
			writeError(writer, err)
			return
		}

		writer.WriteHeader(200)
		_, err = writer.Write(respBytes)

		// we already wrote the header, and cannot write again
		// if any error, write the error object
		if err != nil {
			writeError(writer, err)
		}

	})
}

func writeError(writer http.ResponseWriter, err error) {
	errObj := RespErr{
		Message: err.Error(),
	}
	errObjBytes, err2 := json.Marshal(errObj)
	if err2 != nil {
		fmt.Println(err2)
	}

	_, err2 = writer.Write(errObjBytes)
	if err2 != nil {
		fmt.Println(err2)
	}
}
