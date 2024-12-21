package service

import (
	"fmt"
	msDocumentCLient "nexbit/external/microsoft"
	"nexbit/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DocService interface {
	ParseDoc(ctx *fiber.Ctx, modelID string, reqBody models.FetchDocumentRequest) (*models.FetchParseDocumentApiResponse, error)
}

type docService struct {
	DIClient msDocumentCLient.DIClientClient
}

func NewDocService(diClient msDocumentCLient.DIClientClient) *docService {
	return &docService{
		DIClient: diClient,
	}
}

func (s *docService) ParseDoc(ctx *fiber.Ctx, modelID string, reqBody models.FetchDocumentRequest) (*models.FetchParseDocumentApiResponse, error) {

	uploadDocumentResp, err := s.DIClient.UploadAnalysisDoc(ctx, modelID, reqBody)
	if err != nil {
		// util.WithContext(ctx.Context()).Errorf("[ChatService] Failed to process chat request. err: %v", err)
		return nil, err
	}

	time.Sleep(40 * time.Second)

	operationLocation := uploadDocumentResp["Operation-Location"]

	fetchDocumentResp, err := s.DIClient.FetchParsedDoc(ctx, operationLocation.(string))
	if err != nil {
		// util.WithContext(ctx.Context()).Errorf("[ChatService] Failed to process chat request. err: %v", err)
		return nil, err
	}

	finalResponse, err := buildFetchDocResponse(fetchDocumentResp.AnalyzeResult)
	if err != nil {
		// util.WithContext(ctx.Context()).Errorf("[ChatService] Failed to process chat request. err: %v", err)
		return nil, err
	}

	return finalResponse, nil
}

func buildFetchDocResponse(resp models.AnalyzeResultDocResponse) (*models.FetchParseDocumentApiResponse, error) {

	processedTableResponse, err := ProcessTableData(resp.Tables)
	if err != nil {
		// util.WithContext(ctx.Context()).Errorf("[ChatService] Failed to process chat request. err: %v", err)
		return nil, err
	}

	fetchDocResponse := models.FetchParseDocumentApiResponse{
		Tables:     processedTableResponse,
		Paragraphs: resp.Paragraphs,
	}

	return &fetchDocResponse, nil
}

func ProcessTableData(fetchedTables []models.Table) ([]models.ProcessedTableResponse, error) {

	var TableList []models.ProcessedTableResponse

	for tableIndex, table := range fetchedTables {
		processedTable, err := processTable(&table)
		if err != nil {
			return nil, fmt.Errorf("error processing table %d: %v", tableIndex, err)
		}
		TableList = append(TableList, *processedTable)
	}

	return TableList, nil
}

func processTable(table *models.Table) (*models.ProcessedTableResponse, error) {
	result := &models.ProcessedTableResponse{
		Headers: make([]string, 0),
		Rows:    make([]map[string]string, 0),
		Meta: map[string]interface{}{
			"rowCount":    table.RowCount,
			"columnCount": table.ColumnCount,
		},
	}

	// Extract headers and create header mapping
	headerMap := make(map[int]string)
	for _, cell := range table.Cells {
		if cell.Kind == "columnHeader" {
			result.Headers = append(result.Headers, cell.Content)
			headerMap[cell.ColumnIndex] = cell.Content
		}
	}

	// Process rows
	rowMap := make(map[int]map[string]string)

	// Initialize row maps
	for i := 0; i < table.RowCount; i++ {
		rowMap[i] = make(map[string]string)
	}

	// Fill in cell values
	for _, cell := range table.Cells {
		if cell.Kind == "columnHeader" {
			continue
		}

		// Get header name for this column
		headerName := headerMap[cell.ColumnIndex]
		if headerName == "" {
			headerName = fmt.Sprintf("column_%d", cell.ColumnIndex)
		}

		// Store cell content
		if rowMap[cell.RowIndex] == nil {
			rowMap[cell.RowIndex] = make(map[string]string)
		}
		rowMap[cell.RowIndex][headerName] = cell.Content
	}

	// Convert row map to slice, skipping header row
	for i := 1; i < table.RowCount; i++ {
		if row := rowMap[i]; len(row) > 0 {
			result.Rows = append(result.Rows, row)
		}
	}

	return result, nil
}
