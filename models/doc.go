package models

type FetchDocumentRequest struct {
	Base64Source string `json:"base64Source"`
}

type FetchDocumentIntelligenceDocResponse struct {
	Status              string                   `json:"status"`
	CreatedDateTime     string                   `json:"createdDateTime"`
	LastUpdatedDateTime string                   `json:"lastUpdatedDateTime"`
	AnalyzeResult       AnalyzeResultDocResponse `json:"analyzeResult"`
}

type AnalyzeResultDocResponse struct {
	ApiVersion      string        `json:"apiVersion"`
	ModelId         string        `json:"modelId"`
	StringIndexType string        `json:"stringIndexType"`
	Pages           []interface{} `json:"pages"`
	Tables          []Table       `json:"tables"`
	Paragraphs      []Paragraph   `json:"paragraphs"`
	Styles          []interface{} `json:"styles"`
	ContentFormat   string        `json:"contentFormat"`
	Sections        []interface{} `json:"sections"`
	Figures         []interface{} `json:"figures"`
}

type TableCell struct {
	Kind        string `json:"kind,omitempty"`
	RowIndex    int    `json:"rowIndex"`
	ColumnIndex int    `json:"columnIndex"`
	Content     string `json:"content"`
}

// Table represents the basic table structure from the input
type Table struct {
	RowCount    int         `json:"rowCount"`
	ColumnCount int         `json:"columnCount"`
	Cells       []TableCell `json:"cells"`
}

// ProcessedTable represents a cleaned table with headers and rows
type ProcessedTableResponse struct {
	Headers []string               `json:"headers"`
	Rows    []map[string]string    `json:"rows"`
	Meta    map[string]interface{} `json:"meta"`
}

// ProcessedData contains all processed tables
type ProcessedData struct {
	Tables []ProcessedTableResponse `json:"tables"`
}

// Paragraph represents each text block with its metadata
type Paragraph struct {
	Content         string        `json:"content"`
	BoundingRegions []BoundingBox `json:"boundingRegions"`
	Spans           []Span        `json:"spans"`
	Role            string        `json:"role,omitempty"`
}

// BoundingBox represents the position of text
type BoundingBox struct {
	PageNumber int       `json:"pageNumber"`
	Polygon    []float64 `json:"polygon"`
}

// Span represents text positioning
type Span struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}

// ProcessedData contains all processed tables
type FetchParseDocumentApiResponse struct {
	Tables     []ProcessedTableResponse `json:"tables"`
	Paragraphs []Paragraph              `json:"paragraphs"`
}
