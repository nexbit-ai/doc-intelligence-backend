package ory

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	external "nexbit/external"
	"nexbit/models"
	"os"

	"github.com/gofiber/fiber/v2"
)

// const API_TOKEN = "4cnepI1BCmxxmbC2wcVaD3yE0EjmOqmWvhS9i3gtAHOLic7h7GpqJQQJ99ALACYeBjFXJ3w3AAALACOGr1PC"
const BASE_URL = "https://nexbitdoc.cognitiveservices.azure.com/"
const API_VERSION = "2024-02-29-preview"

type DIClientClient interface {
	UploadAnalysisDoc(ctx *fiber.Ctx, modelID string, reqBody models.FetchDocumentRequest) (map[string]interface{}, error)
	FetchParsedDoc(ctx *fiber.Ctx, operationLocationURL string) (*models.FetchDocumentIntelligenceDocResponse, error)
}

func NewDIClientClient(httpClient *external.HTTPClient) *diClientClient {
	return &diClientClient{
		httpClient: httpClient,
	}
}

type diClientClient struct {
	httpClient *external.HTTPClient
}

func (c *diClientClient) UploadAnalysisDoc(ctx *fiber.Ctx, modelID string, reqBody models.FetchDocumentRequest) (map[string]interface{}, error) {
	url := fmt.Sprintf("%sdocumentintelligence/documentModels/%s:analyze?api-version=%s",
		BASE_URL, modelID, API_VERSION)

	var documentIntelligenceApiKey = os.Getenv("DOCUMENT_INTELLIGENCE_API_KEY")

	log.Printf("[UploadAnalysisDoc] Initiating document analysis with URL: %s", url)

	headers := map[string]string{
		"Content-Type":              "application/json",
		"Ocp-Apim-Subscription-Key": documentIntelligenceApiKey,
	}

	reqBodyByte, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("[UploadAnalysisDoc] failed to marshal request body: %w", err)
	}

	resp, err := c.httpClient.Post(ctx.Context(), url, headers, reqBodyByte)
	if err != nil {
		return nil, fmt.Errorf("[UploadAnalysisDoc] error making POST request: %w", err)
	}

	defer resp.Body.Close() // Ensure response body is closed

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("[UploadAnalysisDoc] unexpected status code: %d", resp.StatusCode)
	}

	log.Println("[UploadAnalysisDoc] Successfully received accepted response from the API")

	// Extract relevant headers with default values if missing
	finalResponse := map[string]interface{}{
		"Operation-Location": getHeader(resp.Header, "Operation-Location"),
		"Date":               getHeader(resp.Header, "Date"),
		"apim-request-id":    getHeader(resp.Header, "apim-request-id"),
	}

	return finalResponse, nil
}

// Helper function to retrieve headers with default value handling
func getHeader(headers map[string][]string, key string) string {
	if values, found := headers[key]; found && len(values) > 0 {
		return values[0]
	}
	return ""
}

func (c *diClientClient) FetchParsedDoc(ctx *fiber.Ctx, operationLocationURL string) (*models.FetchDocumentIntelligenceDocResponse, error) {
	log.Printf("[FetchParsedDoc] Calling fetch doc API with URL: %s", operationLocationURL)
	var documentIntelligenceApiKey = os.Getenv("DOCUMENT_INTELLIGENCE_API_KEY")
	headers := map[string]string{
		"Content-Type":              "application/json",
		"Ocp-Apim-Subscription-Key": documentIntelligenceApiKey,
	}

	resp, err := c.httpClient.Get(ctx.Context(), operationLocationURL, headers)
	if err != nil {
		return nil, fmt.Errorf("[FetchParsedDoc] error fetching data: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[FetchParsedDoc] received non-200 status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[FetchParsedDoc] failed to read response body: %w", err)
	}

	var respBody models.FetchDocumentIntelligenceDocResponse
	if err := json.Unmarshal(respBodyByte, &respBody); err != nil {
		return nil, fmt.Errorf("[FetchParsedDoc] failed to unmarshal response body: %w", err)
	}

	log.Println("[FetchParsedDoc] Successfully fetched and parsed document intelligence response")
	return &respBody, nil
}
