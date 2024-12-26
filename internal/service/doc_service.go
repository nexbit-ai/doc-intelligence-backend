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

	time.Sleep(8 * time.Second)

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

func buildFetchDocResponse(resp models.InvoiceModelAnalyzeResultResponse) (*models.FetchParseDocumentApiResponse, error) {

	processedTableResponse, err := ProcessTableData(resp.Tables)
	if err != nil {
		// util.WithContext(ctx.Context()).Errorf("[ChatService] Failed to process chat request. err: %v", err)
		return nil, err
	}

	fetchDocResponse := models.FetchParseDocumentApiResponse{
		Tables:    processedTableResponse,
		Documents: resp.Documents,
	}

	return &fetchDocResponse, nil
}

func ProcessTableData(fetchedTables []models.Table) ([]models.ProcessedInvoice, error) {

	var TableList []models.ProcessedInvoice

	fmt.Println(len(fetchedTables))
	for tableIndex, table := range fetchedTables {
		processedTable, err := ProcessInvoiceData(table)
		if err != nil {
			return nil, fmt.Errorf("error processing table %d: %v", tableIndex, err)
		}

		if len(processedTable.Headers) != 0 {
			TableList = append(TableList, *processedTable)
		}
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

// determineValueType tries to identify the type of value in a cell
func determineValueType(content string) string {
	if content == "" {
		return "text"
	}

	// Check if it's a percentage
	if len(content) >= 1 && content[len(content)-1] == '%' {
		return "percentage"
	}

	// Check if it contains currency symbols or decimal points
	if len(content) > 0 {
		firstChar := content[0]
		if (firstChar >= '0' && firstChar <= '9') || firstChar == '-' || firstChar == ',' {
			return "amount"
		}
	}

	return "text"
}

// extractHeaderHierarchy processes cells to create a hierarchical header structure
func extractHeaderHierarchy(table models.Table) []models.HeaderColumn {
	headerRows := make(map[int][]models.TableCell)
	maxHeaderRow := -1

	// Group header cells by row
	for _, cell := range table.Cells {
		if cell.Kind == "columnHeader" {
			headerRows[cell.RowIndex] = append(headerRows[cell.RowIndex], cell)
			if cell.RowIndex > maxHeaderRow {
				maxHeaderRow = cell.RowIndex
			}
		}
	}

	// Process header hierarchy
	headers := make([]models.HeaderColumn, 0)
	processedColumns := make(map[int]bool)

	// Process each column in the first row
	for _, cell := range headerRows[0] {
		if processedColumns[cell.ColumnIndex] {
			continue
		}

		header := models.HeaderColumn{
			Title: cell.Content,
			Span:  cell.ColumnSpan,
			Level: 0,
		}

		// Check for child headers
		if cell.RowSpan < maxHeaderRow+1 {
			children := extractChildHeaders(headerRows, cell, 1, maxHeaderRow)
			fmt.Println("reeeeeached here")
			fmt.Println(children)
			header.Children = children
		}

		headers = append(headers, header)
		processedColumns[cell.ColumnIndex] = true
	}

	return headers
}

// extractChildHeaders recursively processes child headers
func extractChildHeaders(headerRows map[int][]models.TableCell, parentCell models.TableCell, currentRow int, maxRow int) []models.HeaderColumn {
	// Debug logging
	fmt.Printf("Processing row %d (max: %d) for parent cell at col %d with span %d\n",
		currentRow, maxRow, parentCell.ColumnIndex, parentCell.ColumnSpan)

	children := make([]models.HeaderColumn, 0)

	// Basic validation
	if currentRow > maxRow {
		fmt.Println("Stopping: currentRow > maxRow")
		return children
	}

	// Validate row exists in map
	if _, exists := headerRows[currentRow]; !exists {
		fmt.Printf("Row %d not found in headerRows\n", currentRow)
		return children
	}

	startCol := parentCell.ColumnIndex
	endCol := startCol + parentCell.ColumnSpan

	// Validate column bounds
	if startCol < 0 || endCol < startCol {
		fmt.Printf("Invalid column bounds: start=%d, end=%d\n", startCol, endCol)
		return children
	}

	col := startCol
	processedCols := make(map[int]bool) // Track processed columns to prevent loops

	for col < endCol {
		if processedCols[col] {
			fmt.Printf("Column %d already processed, advancing\n", col)
			col++
			continue
		}

		fmt.Printf("Looking for cell at column %d\n", col)
		cellFound := false

		for _, cell := range headerRows[currentRow] {
			if cell.ColumnIndex == col {
				fmt.Printf("Found cell at col %d with content: %s\n", col, cell.Content)

				header := models.HeaderColumn{
					Title: cell.Content,
					Span:  cell.ColumnSpan,
					Level: currentRow,
				}

				nextRow := currentRow + cell.RowSpan
				if nextRow <= maxRow && cell.RowSpan < maxRow-currentRow+1 {
					fmt.Printf("Recursing to row %d\n", nextRow)
					childHeaders := extractChildHeaders(headerRows, cell, nextRow, maxRow)
					header.Children = childHeaders
				}

				children = append(children, header)
				processedCols[col] = true
				col += cell.ColumnSpan
				cellFound = true
				break
			}
		}

		if !cellFound {
			fmt.Printf("No cell found at column %d, advancing\n", col)
			processedCols[col] = true
			col++
		}
	}

	fmt.Printf("Returning %d children for row %d\n", len(children), currentRow)
	return children
}

// ProcessInvoiceData converts raw invoice analysis to a structured format
func ProcessInvoiceData(table models.Table) (*models.ProcessedInvoice, error) {

	// Get header structure
	fmt.Println("here1")
	headers := extractHeaderHierarchy(table)

	fmt.Println("here")
	fmt.Println(headers)

	// Initialize processed invoice
	processed := &models.ProcessedInvoice{
		Headers:     headers,
		Rows:        make([]models.InvoiceRow, 0),
		ColumnTypes: make(map[int]string),
	}

	// Create a map to store cells by row index
	rowMap := make(map[int][]models.TableCell)
	headerRowCount := 0
	for _, cell := range table.Cells {
		if cell.Kind == "columnHeader" {
			if cell.RowIndex+1 > headerRowCount {
				headerRowCount = cell.RowIndex + 1
			}
			continue
		}
		rowMap[cell.RowIndex] = append(rowMap[cell.RowIndex], cell)
	}

	// Process data rows
	for rowIdx := headerRowCount; rowIdx < table.RowCount; rowIdx++ {
		cells := rowMap[rowIdx]
		if len(cells) == 0 {
			continue
		}

		row := processDataRow(cells, processed.ColumnTypes)
		processed.Rows = append(processed.Rows, row)
	}

	return processed, nil
}

// processDataRow converts table cells into a structured row
func processDataRow(cells []models.TableCell, columnTypes map[int]string) models.InvoiceRow {
	row := models.InvoiceRow{
		Values:  make(map[string]models.InvoiceValue),
		IsTotal: false,
	}

	for _, cell := range cells {
		// Generate a unique identifier for the column
		colID := generateColumnID(cell.ColumnIndex)

		// Determine value type
		valueType := determineValueType(cell.Content)
		columnTypes[cell.ColumnIndex] = valueType

		// Check if this is a total row
		if cell.Content == "Grand Total" || cell.Content == "Total" {
			row.IsTotal = true
		}

		// Store the value
		row.Values[colID] = models.InvoiceValue{
			Value: cell.Content,
			Type:  valueType,
		}
	}

	return row
}

// generateColumnID creates a unique identifier for a column
func generateColumnID(columnIndex int) string {
	return fmt.Sprintf("col_%d", columnIndex)
}
