package rest_test

import "net/http"

const (
	testUuid   = "03cb753f-5e03-4c97-8e47-625115476c72"
	testNodeId = "uk-sa2341-hnode-v0-a1a0"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}
