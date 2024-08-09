package paginatorclient

import "github.com/valyala/fasthttp"

type HostClient interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}
