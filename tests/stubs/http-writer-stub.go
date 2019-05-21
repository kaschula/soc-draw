package stubs

import "net/http"

type HttpWriterStub struct {
	WrittenData [][]byte
}

func (w *HttpWriterStub) Header() http.Header {
	return http.Header{}
}

func (w *HttpWriterStub) Write(data []byte) (int, error) {
	w.WrittenData = append(w.WrittenData, data)

	return len(data), nil
}

func (w *HttpWriterStub) WriteHeader(statusCode int) {

}
