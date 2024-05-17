package model

type Attribute struct {
	Name  string
	Value string
}

type AttributeGroup struct {
	GroupName  string
	Attributes []Attribute `xml:"Attribute"`
}

type Metadata struct {
	Attributes      []Attribute      `xml:"Attribute,omitempty"`
	AttributeGroups []AttributeGroup `xml:"AttributeGroup,omitempty"`
}

type Request struct {
	RowNumber int    `xml:"-"`
	RefID     string `xml:",attr"`
	FileName  string `xml:",attr"`
	MimeType  string `xml:",attr"`
	DocName   string `xml:",attr,omitempty"`
	//must be pointer since golang is mainly value-based.
	//If not pointer, it creates default Metadata struct with default values for all of its fields
	Metadata *Metadata `xml:",omitempty"`
}
