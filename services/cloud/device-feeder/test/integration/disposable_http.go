package integration

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type DisposableServer struct {
	port          int
	requestsCount int
}

func (d *DisposableServer) WaitForRequest(forHowLong time.Duration) (isRequestReceived bool) {
	ctx, cancel := context.WithTimeout(context.Background(), forHowLong)
	defer cancel()

	select {
	case <-ctx.Done():
		d.log("DisposableServer: timeout")
		return false

	case res := <-d.runServe(d.port):
		d.log("received")
		return res
	}
}
func (d *DisposableServer) log(str ...interface{}) {
	fmt.Println("DisposableServer ", str)
}

func (d *DisposableServer) runServe(port int) chan bool {
	ch := make(chan bool)
	srv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	srv.SetKeepAlivesEnabled(false)
	counter := 0

	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.EqualFold(r.URL.Path, "/testEndpoint") {
			d.log("Received request: %s", r.URL.Path)
			_, err := w.Write([]byte("hello from disposable http server"))
			if err != nil {
				d.log("Error writing response:", err.Error())
			}
			counter++
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

		if counter >= d.requestsCount {
			ch <- true
		}
	})

	go func() {
		d.log("Starting server on port ", port)
		err := srv.ListenAndServe()
		if err != nil {
			d.log("error serving: ", err.Error())
			ch <- false
		}

		ch <- true
	}()

	return ch
}
