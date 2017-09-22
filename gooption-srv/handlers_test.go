

package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gooption/pb"
)


func Test_PriceHandler(t *testing.T) {
	router := gin.Default()
	router.POST("/price", handlerPrice)

	if file, err := os.Open("./testdata/PriceRequest.json"); err == nil {
		defer file.Close()
		httpRequest, _ := http.NewRequest("POST", "/price", file)
		httpRequest.Header.Set("Content-Type", binding.MIMEJSON)

		httpResponse := httptest.NewRecorder()
		router.ServeHTTP(httpResponse, httpRequest)

		response := &pb.PriceResponse{}
		if err = jsonpb.Unmarshal(httpResponse.Body, response); err != nil {
			t.Error(err)
		}
		t.Log(response.String())
	} else {
		t.Error(err)
	}
}

func Test_GreekHandler(t *testing.T) {
	router := gin.Default()
	router.POST("/greek", handlerGreek)

	if file, err := os.Open("./testdata/GreekRequest.json"); err == nil {
		defer file.Close()
		httpRequest, _ := http.NewRequest("POST", "/greek", file)
		httpRequest.Header.Set("Content-Type", binding.MIMEJSON)

		httpResponse := httptest.NewRecorder()
		router.ServeHTTP(httpResponse, httpRequest)

		response := &pb.GreekResponse{}
		if err = jsonpb.Unmarshal(httpResponse.Body, response); err != nil {
			t.Error(err)
		}
		t.Log(response.String())
	} else {
		t.Error(err)
	}
}

func Test_ImpliedVolHandler(t *testing.T) {
	router := gin.Default()
	router.POST("/impliedvol", handlerImpliedVol)

	if file, err := os.Open("./testdata/ImpliedVolRequest.json"); err == nil {
		defer file.Close()
		httpRequest, _ := http.NewRequest("POST", "/impliedvol", file)
		httpRequest.Header.Set("Content-Type", binding.MIMEJSON)

		httpResponse := httptest.NewRecorder()
		router.ServeHTTP(httpResponse, httpRequest)

		response := &pb.ImpliedVolResponse{}
		if err = jsonpb.Unmarshal(httpResponse.Body, response); err != nil {
			t.Error(err)
		}
		t.Log(response.String())
	} else {
		t.Error(err)
	}
}

