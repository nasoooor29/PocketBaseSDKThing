package main

// Field related Structs

type PbTextFieldOptions struct {
	Min     int    `json:"min"` // null by default or there is not value
	Max     int    `json:"max"` // null by default or there is not value
	Pattern string `json:"pattern"`
}
type PbFileFieldOptions struct {
	MimeTypes []string    `json:"mimeTypes"`
	Thumbs    interface{} `json:"thumbs"`
	MaxSelect int         `json:"maxSelect"`
	MaxSize   int         `json:"maxSize"`
	Protected bool        `json:"protected"`
}
type PbRelationFieldOptions struct {
	CollectionID string `json:"collectionId"` // the id of the collection,
	// i think we should save the collection names into a map for easy for the ids access
	CascadeDelete bool        `json:"cascadeDelete"` // it's false by default
	MinSelect     interface{} `json:"minSelect"`     // null by default
	MaxSelect     int         `json:"maxSelect"`     // i think 1 for one to one
	DisplayFields interface{} `json:"displayFields"` // null by default
}

type PbEditorFieldOptions struct {
	ConvertUrls bool `json:"convertUrls"` // false by default
}
type PbNumberFieldOptions struct {
	Min       int  `json:"min"`       // null by default
	Max       int  `json:"max"`       // null by default
	NoDecimal bool `json:"noDecimal"` // false by default
}

type PbBoolFieldOptions struct{}

type PbDomainFieldOptions struct { // it's the same for url type and email
	ExceptDomains []string `json:"exceptDomains"` // should be null if it's not set
	OnlyDomains   []string `json:"onlyDomains"`   // should be null if it's not set
}

type PbDateFieldOptions struct {
	Min string `json:"min"` // null by default, if set it should be in utc time "2024-10-26 12:00:00.000Z"
	Max string `json:"max"` // null by default, if set it should be in utc time "2024-10-26 12:00:00.000Z"
}
type PbSelectFieldOptions struct {
	MaxSelect int      `json:"maxSelect"` // number of max selections, 1 by default
	Values    []string `json:"values"`    // the available values, null by default
}

type PbJsonFieldOptions struct {
	MaxSize int `json:"maxSize"` // 2000000 by default
}

type PbFieldType interface {
	PbTextFieldOptions |
		PbFileFieldOptions |
		PbRelationFieldOptions |
		PbEditorFieldOptions |
		PbNumberFieldOptions |
		PbBoolFieldOptions |
		PbDomainFieldOptions |
		PbDateFieldOptions |
		PbSelectFieldOptions |
		PbJsonFieldOptions |
		any // this is just to shut the error fix later
}

type PbField[fieldType PbFieldType] struct { // it's the fields inside the struct
	*FieldData `json:"fieldData"` // it's the same for all types
	Options   fieldType `json:"options"`   // it would be different for each type
}

type FieldData struct {
	System      bool   `json:"system"`      // same as above
	ID          string `json:"id"`          // auto generated
	Name        string `json:"name"`        // field name
	Type        string `json:"type"`        // the available types on pocketbase so we should make a map and custom types for this
	Required    bool   `json:"required"`    // if the field is required when creating new object it's off by default
	Presentable bool   `json:"presentable"` // it's false by default
	Unique      bool   `json:"unique"`      // it's false by default
}

// Collection related Structs

type PbAuthCollectionOptions struct {
	AllowEmailAuth     bool        `json:"allowEmailAuth"`
	AllowOAuth2Auth    bool        `json:"allowOAuth2Auth"`
	AllowUsernameAuth  bool        `json:"allowUsernameAuth"`
	ExceptEmailDomains interface{} `json:"exceptEmailDomains"`
	ManageRule         interface{} `json:"manageRule"`
	MinPasswordLength  int         `json:"minPasswordLength"`
	OnlyEmailDomains   interface{} `json:"onlyEmailDomains"`
	OnlyVerified       bool        `json:"onlyVerified"`
	RequireEmail       bool        `json:"requireEmail"`
}

type PbViewCollectionOptions struct {
	Query string `json:"query"`
}

type PbBaseCollectionOptions struct{}

type CollectionType interface {
	PbAuthCollectionOptions | PbBaseCollectionOptions | PbViewCollectionOptions
}

type PocketBaseCollection[collType CollectionType] struct {
	ID         string        `json:"id"`     // auto generated
	Name       string        `json:"name"`   // struct Name
	Type       string        `json:"type"`   // base, auth, view
	System     bool          `json:"system"` // tbh, I don't know but all the time is false
	Schema     []any         `json:"schema"`
	Indexes    []interface{} `json:"indexes"`    // for later
	ListRule   string        `json:"listRule"`   // empty by default
	ViewRule   string        `json:"viewRule"`   // empty by default
	CreateRule string        `json:"createRule"` // empty by default
	UpdateRule string        `json:"updateRule"` // empty by default
	DeleteRule string        `json:"deleteRule"` // empty by default
	Options    collType      `json:"options"`    // it would be different for each type
}
