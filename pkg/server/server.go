package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Big-Kotik/ivt-pull-api/pkg/api"
)

type PullServer struct {
	http.Client
	Logger *log.Logger
	api.UnimplementedPullerServer
}

func (p *PullServer) PullResources(req *api.HttpRequests, respStream api.Puller_PullResourcesServer) error {
	// p.Logger.Printf("start handling requests")

	for _, req := range req.GetRequests() {
		fmt.Printf("pull: %s\n", req.Url)
		resp, err := p.sendRequset(req)

		if err != nil {
			p.Logger.Print(fmt.Printf("can't send request in PullResources - %s", err.Error()))

			fmt.Println("error")
			return err
		}

		fmt.Printf("resp UUid: %v\n", req.Uuid)
		resp.Uuid = req.Uuid
		err = respStream.Send(resp)

		if err != nil {
			p.Logger.Print(fmt.Printf("can't write resp to stream - %s", err.Error()))

			fmt.Println("another error")
			return err
		}

		fmt.Printf("pulled: %s\n", req.Url)
	}

	// fmt.Println("exit")
	return nil
}

func (p *PullServer) mustEmbedUnimplementedPullerServer() {

}

func newHttpRequest(r *api.HttpRequests_HttpRequest) (*http.Request, error) {
	req, err := http.NewRequest(r.GetMethod(), r.GetUrl(), bytes.NewReader(r.GetBody()))

	if err != nil {
		return nil, err
	}

	for k, v := range r.GetHeaders() {
		for _, key := range v.Keys {
			req.Header.Set(k, key)
		}
	}

	return req, nil
}

func (p *PullServer) sendRequset(r *api.HttpRequests_HttpRequest) (*api.HttpResponse, error) {
	req, err := newHttpRequest(r)

	if err != nil {
		p.Logger.Printf("Bad request in sendRequest - %s", err.Error())

		return nil, fmt.Errorf("bad request in sendRequest - %w", err)
	}

	resp, err := p.Client.Do(req)

	if err != nil {
		p.Logger.Printf("Can't send request in sendRequest - %s", err.Error())

		return nil, fmt.Errorf("can't send request in sendRequest - %w", err)
	}

	return httpToGrpc(resp)
}

func httpToGrpc(resp *http.Response) (*api.HttpResponse, error) {
	grpcResp := api.HttpResponse{
		StatusCode: int32(resp.StatusCode),
		ProtoMajor: int32(resp.ProtoMajor),
		ProtoMinor: int32(resp.ProtoMinor),
	}
	grpcResp.Header = make(map[string]*api.Header)

	for k, v := range resp.Header {
		grpcResp.Header[k] = &api.Header{Keys: v}
	}

	body, err := io.ReadAll(resp.Body)
	grpcResp.Body = body

	return &grpcResp, err
}
