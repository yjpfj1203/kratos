package http

import (
	"net/http"
)

// DefaultResponseEncoder is default response encoder.
func DefaultResponseEncoder(res http.ResponseWriter, req *http.Request, v interface{}) error {
	contentType, codec, err := responseCodec(req)
	if err != nil {
		return err
	}
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	res.Header().Set("content-type", contentType)
	res.Write(data)
	return nil
}

// DefaultErrorEncoder is default errors encoder.
func DefaultErrorEncoder(res http.ResponseWriter, req *http.Request, err error) {
	code, se := StatusError(err)
	contentType, codec, err := responseCodec(req)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := codec.Marshal(se)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", contentType)
	res.WriteHeader(code)
	res.Write(data)
}
