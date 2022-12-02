package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type sampleResponse struct {
	Payload string `json:"payload"`
}

// emptyHandler is a simple empty test handler
func emptyHandler() func(http.ResponseWriter, *http.Request) {
	hfn := func(r *Request) (ResponseData, error) {
		return nil, nil
	}
	return Handle(hfn)
}

// stringResponseHandler this test handler response with string
func stringResponseHandler(respData string) func(http.ResponseWriter, *http.Request) {
	hfn := func(r *Request) (ResponseData, error) {
		_, _ = r.W.Write([]byte(respData))
		return nil, nil
	}

	return Handle(hfn)
}

// jsonResponseHandler response with json data
func jsonResponseHandler(respData sampleResponse) func(http.ResponseWriter, *http.Request) {
	hfn := func(r *Request) (ResponseData, error) {
		return respData, nil
	}

	return Handle(hfn)
}

func TestHandle_EmptyHandler(t *testing.T) {
	resp := setupRecord(emptyHandler(), http.MethodGet, nil)
	defer resp.Body.Close()

	err := checkBody(resp.Body, "")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestHandle_StringResponse(t *testing.T) {
	sampleResponseString := "myPayload"
	resp := setupRecord(stringResponseHandler(sampleResponseString), http.MethodPost, nil)
	defer resp.Body.Close()

	wantCode := 200
	if wantCode != resp.StatusCode {
		t.Fatalf("unexpected response code, want: %d, got: %d\n", wantCode, resp.StatusCode)
	}

	err := checkBody(resp.Body, sampleResponseString)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestJsonParser_JsonResponse(t *testing.T) {
	sampleResp := sampleResponse{"myPayload"}
	resp := setupRecord(jsonResponseHandler(sampleResp), http.MethodPost, nil)
	defer resp.Body.Close()

	wantCode := 200
	if wantCode != resp.StatusCode {
		t.Fatalf("unexpected response code, want: %d, got: %d\n", wantCode, resp.StatusCode)
	}

	sampleRespByte, _ := json.Marshal(sampleResp)

	err := checkBody(resp.Body, string(sampleRespByte))
	if err != nil {
		t.Error(err.Error())
	}
}

func setupRecord(handlerFn func(http.ResponseWriter, *http.Request), method string, body io.Reader) *http.Response {
	req := httptest.NewRequest(method, "/", body)
	recorder := httptest.NewRecorder()
	handlerFn(recorder, req)
	return recorder.Result()
}

func checkBody(body io.ReadCloser, expected string) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	if string(data) != expected {
		return errors.New(fmt.Sprintf("unexpected result: '%s', got: '%s'", expected, data))
	}
	return nil
}
