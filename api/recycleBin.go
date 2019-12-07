package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/koltyakov/gosip"
)

// RecycleBin ...
type RecycleBin struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// RecycleBinItem ...
type RecycleBinItem struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// RecycledItem ...
type RecycledItem struct {
	AuthorEmail               string    `json:"AuthorEmail"`
	AuthorName                string    `json:"AuthorName"`
	DeletedByEmail            string    `json:"DeletedByEmail"`
	DeletedByName             string    `json:"DeletedByName"`
	DeletedDate               time.Time `json:"DeletedDate"`
	DeletedDateLocalFormatted string    `json:"DeletedDateLocalFormatted"`
	DirName                   string    `json:"DirName"`
	ID                        string    `json:"Id"`
	ItemState                 int       `json:"ItemState"`
	ItemType                  int       `json:"ItemType"`
	LeafName                  string    `json:"LeafName"`
	Size                      int       `json:"Size"`
	Title                     string    `json:"Title"`
	LeafNamePath              struct {
		DecodedURL string `json:"DecodedUrl"`
	} `json:"LeafNamePath"`
	DirNamePath struct {
		DecodedURL string `json:"DecodedUrl"`
	} `json:"DirNamePath"`
}

// RecycleBinResp ...
type RecycleBinResp []byte

// RecycleBinItemResp ...
type RecycleBinItemResp []byte

// NewRecycleBin ...
func NewRecycleBin(client *gosip.SPClient, endpoint string, config *RequestConfig) *RecycleBin {
	return &RecycleBin{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// NewRecycleBinItem ...
func NewRecycleBinItem(client *gosip.SPClient, endpoint string, config *RequestConfig) *RecycleBinItem {
	return &RecycleBinItem{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (recycleBin *RecycleBin) ToURL() string {
	apiURL, _ := url.Parse(recycleBin.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range recycleBin.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (recycleBin *RecycleBin) Conf(config *RequestConfig) *RecycleBin {
	recycleBin.config = config
	return recycleBin
}

// Select ...
func (recycleBin *RecycleBin) Select(oDataSelect string) *RecycleBin {
	if recycleBin.modifiers == nil {
		recycleBin.modifiers = make(map[string]string)
	}
	recycleBin.modifiers["$select"] = oDataSelect
	return recycleBin
}

// Expand ...
func (recycleBin *RecycleBin) Expand(oDataExpand string) *RecycleBin {
	if recycleBin.modifiers == nil {
		recycleBin.modifiers = make(map[string]string)
	}
	recycleBin.modifiers["$expand"] = oDataExpand
	return recycleBin
}

// Filter ...
func (recycleBin *RecycleBin) Filter(oDataFilter string) *RecycleBin {
	if recycleBin.modifiers == nil {
		recycleBin.modifiers = make(map[string]string)
	}
	recycleBin.modifiers["$filter"] = oDataFilter
	return recycleBin
}

// Top ...
func (recycleBin *RecycleBin) Top(oDataTop int) *RecycleBin {
	if recycleBin.modifiers == nil {
		recycleBin.modifiers = make(map[string]string)
	}
	recycleBin.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return recycleBin
}

// OrderBy ...
func (recycleBin *RecycleBin) OrderBy(oDataOrderBy string, ascending bool) *RecycleBin {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if recycleBin.modifiers == nil {
		recycleBin.modifiers = make(map[string]string)
	}
	if recycleBin.modifiers["$orderby"] != "" {
		recycleBin.modifiers["$orderby"] += ","
	}
	recycleBin.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return recycleBin
}

// Get ...
func (recycleBin *RecycleBin) Get() (RecycleBinResp, error) {
	sp := NewHTTPClient(recycleBin.client)
	return sp.Get(recycleBin.ToURL(), getConfHeaders(recycleBin.config))
}

// GetByID ...
func (recycleBin *RecycleBin) GetByID(itemID string) *RecycleBinItem {
	return NewRecycleBinItem(
		recycleBin.client,
		fmt.Sprintf("%s('%s')", recycleBin.endpoint, itemID),
		recycleBin.config,
	)
}

// Get ...
func (recycleBinItem *RecycleBinItem) Get() (RecycleBinItemResp, error) {
	sp := NewHTTPClient(recycleBinItem.client)
	return sp.Get(recycleBinItem.endpoint, getConfHeaders(recycleBinItem.config))
}

// Restore ...
func (recycleBinItem *RecycleBinItem) Restore() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/Restore()", recycleBinItem.endpoint)
	sp := NewHTTPClient(recycleBinItem.client)
	return sp.Post(endpoint, nil, getConfHeaders(recycleBinItem.config))
}

/* Response helpers */

// Data : to get typed data
func (recycleBinResp *RecycleBinResp) Data() []RecycleBinItemResp {
	collection, _ := parseODataCollection(*recycleBinResp)
	items := []RecycleBinItemResp{}
	for _, item := range collection {
		items = append(items, RecycleBinItemResp(item))
	}
	return items
}

// Unmarshal : to unmarshal to custom object
func (recycleBinResp *RecycleBinResp) Unmarshal(obj interface{}) error {
	data, _ := parseODataCollectionPlain(*recycleBinResp)
	return json.Unmarshal(data, obj)
}

// Data : to get typed data
func (recycleBinItemResp *RecycleBinItemResp) Data() *RecycledItem {
	data := parseODataItem(*recycleBinItemResp)
	res := &RecycledItem{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (recycleBinItemResp *RecycleBinItemResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*recycleBinItemResp)
	data = normalizeMultiLookups(data)
	return json.Unmarshal(data, obj)
}