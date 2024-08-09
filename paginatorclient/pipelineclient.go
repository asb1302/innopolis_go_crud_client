package paginatorclient

import "github.com/valyala/fasthttp"

type PipelineClient interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}
