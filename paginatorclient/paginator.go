package paginatorclient

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/valyala/fasthttp"
)

type PaginatorClient struct {
	hostClient     HostClient
	pipelineClient PipelineClient
	config         PaginatorClientConfig
}

func NewPaginatorClient(hostClient HostClient, pipelineClient PipelineClient, config PaginatorClientConfig) *PaginatorClient {
	return &PaginatorClient{
		hostClient:     hostClient,
		pipelineClient: pipelineClient,
		config:         config,
	}
}

func (pc *PaginatorClient) getTotalCount() (int, error) {
	if pc.hostClient == nil {
		return 0, fmt.Errorf("hostClient is not initialized")
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fmt.Sprintf("%s/count", pc.config.CrudClientURL))
	req.Header.Set("Authorization", pc.config.AuthToken)

	err := pc.hostClient.Do(req, resp)
	if err != nil {
		log.Printf("Request to /count failed: %v", err)

		return 0, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		log.Printf("Unexpected status code from /count: %d", resp.StatusCode())

		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	count, err := strconv.Atoi(string(resp.Body()))
	if err != nil {
		log.Printf("Error parsing count from response: %v", err)

		return 0, err
	}

	return count, nil
}

func (pc *PaginatorClient) GetAllRecipes() ([]map[string]interface{}, error) {
	// Получаем общее количество рецептов
	totalCount, err := pc.getTotalCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	log.Printf("Total recipes count: %d", totalCount)

	// Пагинируем запросы для получения всех рецептов
	recipes := make([]map[string]interface{}, 0)
	pageSize := 10
	totalPages := (totalCount + pageSize - 1) / pageSize // Высчитываем количество страниц

	// Создаем массив для хранения всех ответов
	responses := make([]*fasthttp.Response, totalPages)

	// Отправляем все запросы без ожидания ответов
	for page := 1; page <= totalPages; page++ {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()

		uri := fmt.Sprintf("%s/?page=%d&limit=%d", pc.config.CrudClientURL, page, pageSize)
		req.SetRequestURI(uri)
		req.Header.Set("Authorization", pc.config.AuthToken)
		req.Header.SetMethod(fasthttp.MethodGet)

		// Отправляем запрос и сохраняем ответ для обработки позже
		err := pc.pipelineClient.Do(req, resp)
		if err != nil {
			log.Printf("Request error: %v", err)

			return nil, err
		}

		responses[page-1] = resp
		fasthttp.ReleaseRequest(req)
	}

	// Обрабатываем все полученные ответы
	for page, resp := range responses {
		if resp.StatusCode() != fasthttp.StatusOK {
			log.Printf("Unexpected status code: %d for page %d. Response body: %s", resp.StatusCode(), page+1, resp.Body())

			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		var pageRecipes []map[string]interface{}
		err := json.Unmarshal(resp.Body(), &pageRecipes)
		if err != nil {
			log.Printf("Error unmarshaling response body: %v", err)

			return nil, err
		}

		log.Printf("Successfully received %d recipes for page %d", len(pageRecipes), page+1)
		recipes = append(recipes, pageRecipes...)

		fasthttp.ReleaseResponse(resp)
	}

	return recipes, nil
}
