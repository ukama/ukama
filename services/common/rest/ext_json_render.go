package rest

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

// Extended json renderer.
// It's same as default renderer with exception that it rendeers protobuf messages without ommitin andy filds.
// Default renderer omits fields that have default values

type ExtJson struct {
	Data   interface{}
	Indent bool
}

// Render (JSON) writes data with custom ContentType.
func (r ExtJson) Render(w http.ResponseWriter) (err error) {
	if err = r.writeJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

// writeContentType (JSON) writes JSON ContentType.
func (r ExtJson) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (r ExtJson) writeJSON(w http.ResponseWriter, obj interface{}) (err error) {
	writeContentType(w, jsonContentType)
	var jsonBytes []byte
	if pr, ok := obj.(proto.Message); ok {
		jsonBytes, err = protojson.MarshalOptions{EmitUnpopulated: true, Multiline: r.Indent}.Marshal(pr)
	} else {
		if r.Indent {
			jsonBytes, err = json.MarshalIndent(obj, "", "  ")
		} else {
			jsonBytes, err = json.Marshal(obj)
		}
	}

	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
