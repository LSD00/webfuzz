package reqgen

import (
	"bufio"
	"fmt"

	"github.com/valyala/fasthttp"
)

type Worker struct {
	domain     string
	tlsEnabled bool
	req        *bufio.Reader
}

type Resp struct {
	BodyData    string
	Len, Status int
}

func NewWorker(tlsEnabled bool, req *bufio.Reader, domain string) *Worker {
	return &Worker{
		tlsEnabled: tlsEnabled,
		req:        req,
		domain:     domain,
	}
}

func (w *Worker) MakeRequest() (Resp, error) {
	request := fasthttp.AcquireRequest()
	err := request.Read(w.req)
	if err != nil {
		return Resp{}, err
	}
	if w.tlsEnabled {
		request.SetRequestURI(fmt.Sprintf("https://%s%s", w.domain, string(request.RequestURI())))
	} else {
		request.SetRequestURI(fmt.Sprintf("http://%s%s", w.domain, string(request.RequestURI())))
	}

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(request, resp)

	bodyBytes, err := resp.BodyUncompressed()
	if err != nil {
		return Resp{}, err
	}
	return Resp{
		BodyData: string(bodyBytes),
		Len:      len(string(bodyBytes)),
		Status:   resp.StatusCode(),
	}, nil
}
