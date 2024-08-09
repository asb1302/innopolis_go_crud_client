package main

import (
	"fmt"
	"github.com/asb1302/innopolis_go_crud_client/paginatorclient"
	"log"
	"net/url"

	"github.com/valyala/fasthttp"
)

func main() {
	paginatorclient.InitConfig()
	cfg := paginatorclient.GetConfig()

	parsedURL, err := url.Parse(cfg.CrudClientURL)
	if err != nil {
		log.Fatalf("Invalid URL in config: %v", err)
	}

	host := parsedURL.Host
	if parsedURL.Port() == "" {
		if parsedURL.Scheme == "https" {
			host = parsedURL.Host + ":443"
		} else {
			host = parsedURL.Host + ":80"
		}
	}

	hostClient := &fasthttp.HostClient{
		Addr:  host,
		IsTLS: parsedURL.Scheme == "https",
	}

	pipelineClient := &fasthttp.PipelineClient{
		Addr:  host,
		IsTLS: parsedURL.Scheme == "https",
	}

	paginatorClient := paginatorclient.NewPaginatorClient(hostClient, pipelineClient, *cfg)

	recipes, err := paginatorClient.GetAllRecipes()
	if err != nil {
		log.Fatalf("Failed to get all recipes: %v", err)
	}

	for _, recipe := range recipes {
		fmt.Printf("Recipe: %v\n", recipe)
	}
}
