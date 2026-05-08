package gen

import (
	"encoding/json"
	"google.golang.org/grpc/encoding"
)

const codecName = "json"

type jsonCodec struct{}

func (jsonCodec) Marshal(v interface{}) ([]byte, error)      { return json.Marshal(v) }
func (jsonCodec) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }
func (jsonCodec) Name() string                               { return codecName }
func init()                                                  { encoding.RegisterCodec(jsonCodec{}) }
func ForceJSONCodec() jsonCodec                              { return jsonCodec{} }
