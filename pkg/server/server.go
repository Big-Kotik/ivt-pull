package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Big-Kotik/ivt-pull-api/pkg/api"
)

type PullServer struct {
	http.Client
	Logger *log.Logger
	api.UnimplementedPullerServer
}

// TODO: resources
// TODO: bytes body HttpRequestWrapper
func (p *PullServer) PullResource(req *api.HttpRequestsWrapper, respStream api.Puller_PullResourceServer) error {
	for _, req := range req.GetRequests() {
		resp, err := p.sendRequset(req)

		// TODO: как понять какой запрос дошел, а какой нет
		if err != nil {
			p.Logger.Print(fmt.Printf("can't send request in PullResource - %s", err.Error()))
			continue
		}

		err = respStream.Send(resp)

		if err != nil {
			p.Logger.Print(fmt.Printf("can't write resp to stream - %s", err.Error()))
		}
	}

	return nil
}

func (p *PullServer) mustEmbedUnimplementedPullerServer() {

}

func (p *PullServer) sendRequset(r *api.HttpRequestsWrapper_Request) (*api.Response, error) {
	// TODO: too much logic in one place
	req, err := http.NewRequest(r.GetMethod(), r.GetUrl(), strings.NewReader(r.GetBody()))

	if err != nil {
		p.Logger.Print(fmt.Printf("can' create req in sendRequest, method: %s url: %s",
			r.GetMethod(), r.GetUrl()))
		return nil, err
	}

	for k, v := range r.GetHeaders() {
		for _, key := range v.Keys {
			req.Header.Set(k, key)
		}
	}

	resp, err := p.Client.Do(req)

	if err != nil {
		return nil, err
	}

	return respToGrpcResponse(resp)
}

func respToGrpcResponse(resp *http.Response) (*api.Response, error) {
	grpcResp := api.Response{}
	grpcResp.StatusCode = int32(resp.StatusCode)
	grpcResp.ProtoMajor = int32(resp.ProtoMajor)
	grpcResp.ProtoMinor = int32(resp.ProtoMinor)
	grpcResp.Header = make(map[string]*api.Header)

	for k, v := range resp.Header {
		grpcResp.Header[k] = &api.Header{Keys: v}
	}

	// TODO: send by parts?
	body, err := io.ReadAll(resp.Body)
	grpcResp.Body = body

	return &grpcResp, err
}
