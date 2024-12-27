package models

type FetchDocumentRequest struct {
	Base64Source string `json:"base64Source"`
}

type FetchDocumentIntelligenceDocResponse struct {
	Status              string                            `json:"status"`
	CreatedDateTime     string                            `json:"createdDateTime"`
	LastUpdatedDateTime string                            `json:"lastUpdatedDateTime"`
	AnalyzeResult       InvoiceModelAnalyzeResultResponse `json:"analyzeResult"`
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

type InvoiceModelAnalyzeResultResponse struct {
	ApiVersion      string        `json:"apiVersion"`
	ModelId         string        `json:"modelId"`
	StringIndexType string        `json:"stringIndexType"`
	Pages           []interface{} `json:"pages"`
	Tables          []Table       `json:"tables"`
	Content         string        `json:"content"`
	Documents       []Document    `json:"documents"`
	Styles          []interface{} `json:"styles"`
	ContentFormat   string        `json:"contentFormat"`
	Sections        []interface{} `json:"sections"`
	Figures         []interface{} `json:"figures"`
}

type TableCell struct {
	Kind        string `json:"kind,omitempty"`
	RowIndex    int    `json:"rowIndex"`
	ColumnIndex int    `json:"columnIndex"`
	RowSpan     int    `json:"rowSpan,omitempty"`
	ColumnSpan  int    `json:"columnSpan,omitempty"`
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
	Tables    []ProcessedInvoice `json:"tables"`
	Documents []Document         `json:"documents"`
}

type HeaderColumn struct {
	Title    string         `json:"title"`
	Children []HeaderColumn `json:"children,omitempty"`
	Span     int            `json:"span"`
	Level    int            `json:"level"`
}

// InvoiceValue represents a generic value in the invoice
type InvoiceValue struct {
	Value string `json:"value"`
	Type  string `json:"type"` // amount, percentage, text, code
}

// InvoiceRow represents a dynamic row of invoice data
type InvoiceRow struct {
	Values  map[string]InvoiceValue `json:"values"` // key is column identifier
	IsTotal bool                    `json:"isTotal"`
}

type ProcessedInvoice struct {
	Headers     []HeaderColumn `json:"headers"`
	Rows        []InvoiceRow   `json:"rows"`
	ColumnTypes map[int]string `json:"columnTypes"` // Maps column index to data type
}

type InvoiceAnalysis struct {
	Tables []Table `json:"tables"`
}

// Document represents the main invoice document
type Document struct {
	DocType         string           `json:"docType"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Fields          Fields           `json:"fields"`
	Confidence      float64          `json:"confidence"`
	Spans           []Span           `json:"spans"`
}

// BoundingRegion represents the physical location of elements in the document
type BoundingRegion struct {
	PageNumber int       `json:"pageNumber"`
	Polygon    []float64 `json:"polygon"`
}

// Fields contains all the invoice fields
type Fields struct {
	InvoiceDate              Field `json:"InvoiceDate"`
	InvoiceId                Field `json:"InvoiceId"`
	InvoiceTotal             Field `json:"InvoiceTotal"`
	Items                    Field `json:"Items"`      // Changed this to Field type
	TaxDetails               Field `json:"TaxDetails"` // Changed this as well
	VendorAddress            Field `json:"VendorAddress"`
	VendorAddressRecipient   Field `json:"VendorAddressRecipient"`
	VendorName               Field `json:"VendorName"`
	VendorTaxId              Field `json:"VendorTaxId"`
	BillingAddress           Field `json:"BillingAddress"`
	BillingAddressRecipient  Field `json:"BillingAddressRecipient"`
	CustomerName             Field `json:"CustomerName"`
	CustomerTaxId            Field `json:"CustomerTaxId"`
	PaymentTerm              Field `json:"PaymentTerm"`
	PurchaseOrder            Field `json:"PurchaseOrder"`
	ShippingAddress          Field `json:"ShippingAddress"`
	ShippingAddressRecipient Field `json:"ShippingAddressRecipient"`
	SubTotal                 Field `json:"SubTotal"`
	TotalTax                 Field `json:"TotalTax"`
}

// Field represents a basic document field
type Field struct {
	Type            string           `json:"type"`
	ValueString     string           `json:"valueString,omitempty"`
	ValueDate       string           `json:"valueDate,omitempty"`
	ValueCurrency   *Currency        `json:"valueCurrency,omitempty"`
	ValueAddress    *Address         `json:"valueAddress,omitempty"`
	ValueArray      []ItemField      `json:"valueArray,omitempty"` // Added this for arrays
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Confidence      float64          `json:"confidence"`
	Spans           []Span           `json:"spans"`
}

// ItemField represents an item in the invoice
type ItemField struct {
	Type            string           `json:"type"`
	ValueObject     ItemObject       `json:"valueObject"`
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Confidence      float64          `json:"confidence"`
	Spans           []Span           `json:"spans"`
}

// ItemObject represents the details of an invoice item
type ItemObject struct {
	Amount      Field `json:"Amount"`
	Description Field `json:"Description"`
	ProductCode Field `json:"ProductCode,omitempty"`
	Tax         Field `json:"Tax,omitempty"`
	TaxRate     Field `json:"TaxRate"`
}

// Currency represents a monetary value with its currency code
type Currency struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

// Address represents a structured address
type Address struct {
	HouseNumber   string `json:"houseNumber"`
	Road          string `json:"road"`
	PostalCode    string `json:"postalCode"`
	City          string `json:"city"`
	State         string `json:"state"`
	StreetAddress string `json:"streetAddress"`
	StateDistrict string `json:"stateDistrict"`
}
