package model

import "encoding/xml"

type Tag struct {
	Name  string
	Value string
}

type TagGroup struct {
	GroupName string
	Tags      []Tag  `xml:"Tag"`
	groupId   string `xml:"-"`
}

type Metadata struct {
	Tags      []Tag      `xml:"Tag,omitempty"`
	TagGroups []TagGroup `xml:"TagGroup,omitempty"`
}

func (md *Metadata) AddTagOrGroupTag(groupId string, groupName string, tagName string, tagValue string) {
	if groupName == "" {
		md.addTag(tagName, tagValue)
	} else {
		md.addGroupTag(groupId, groupName, tagName, tagValue)
	}
}

func (md *Metadata) addTag(tagName string, tagValue string) {
	if md.Tags == nil {
		md.Tags = make([]Tag, 0)
	}
	md.Tags = append(md.Tags, Tag{Name: tagName, Value: tagValue})
}

func (md *Metadata) addGroupTag(groupId string, groupName string, tagName string, tagValue string) {
	tagGroup := md.locateOrCreateTagGroup(groupId, groupName)
	if tagGroup.Tags == nil {
		tagGroup.Tags = []Tag{Tag{Name: tagName, Value: tagValue}}
	} else {
		tagGroup.Tags = append(tagGroup.Tags, Tag{Name: tagName, Value: tagValue})
	}
}

func (md *Metadata) locateOrCreateTagGroup(groupId string, groupName string) *TagGroup {
	var tagGroup *TagGroup = nil
	if md.TagGroups == nil {
		md.TagGroups = make([]TagGroup, 0)
	} else {
		for i := range md.TagGroups {
			if groupId == "" && md.TagGroups[i].GroupName == groupName && md.TagGroups[i].groupId == "" {
				tagGroup = &md.TagGroups[i]
				break
			} else if groupId != "" && md.TagGroups[i].groupId == groupId {
				tagGroup = &md.TagGroups[i]
				break
			}
		}
	}
	if tagGroup == nil {
		md.TagGroups = append(md.TagGroups, TagGroup{
			GroupName: groupName,
			groupId:   groupId,
		})
		tagGroup = &md.TagGroups[len(md.TagGroups)-1]
	}
	return tagGroup
}

type Request struct {
	RowNumber int    `xml:"-"`
	ID        string `xml:",attr"`
	FileName  string `xml:",attr"`
	MimeType  string `xml:",attr"`
	DocName   string `xml:",attr,omitempty"`
	//must be pointer since golang is mainly value-based.
	//If not pointer, it creates default Metadata struct with default values for all of its fields
	Metadata *Metadata `xml:",omitempty"`
}

func (r *Request) GetTagValue(name string) string {
	if r.Metadata == nil || r.Metadata.Tags == nil {
		return ""
	}
	return doGetTagValue(&r.Metadata.Tags, name)
}

func doGetTagValue(tags *[]Tag, name string) string {
	for _, tag := range *tags {
		if tag.Name == name {
			return tag.Value
		}
	}
	return ""
}

func (r *Request) GetTagGroupValue(groupId string, groupName string, tagName string) string {
	if r.Metadata == nil || r.Metadata.TagGroups == nil {
		return ""
	}
	tagGroup := r.Metadata.locateOrCreateTagGroup(groupId, groupName)
	if tagGroup == nil || tagGroup.Tags == nil {
		return ""
	}
	return doGetTagValue(&tagGroup.Tags, tagName)
}

type PkgHeader struct {
	SubmissionDate string
	SubmissionTime string
	Source         string
}

type PkgTrailer struct {
	RequestCount int16
}

type Pkg struct {
	XMLName  xml.Name `xml:"Package"`
	ID       string   `xml:"ID,attr"`
	Header   PkgHeader
	Requests []Request `xml:"Requests>Request"`
	Trailer  PkgTrailer
}

type Report struct {
	XMLName   xml.Name `xml:"REPORT"`
	ID        string   `xml:"ID,attr"`
	Header    ReportHeader
	Documents []ReportDocument `xml:"Documents>Document"`
	Trailer   ReportTrailer
}

type ReportHeader struct {
	SubmissionDate     string
	SubmissionTime     string
	RequestApplication string
	PackageName        string
	ContentType        string
	ProcessingDuration string
	ProcessingDate     string
	ProcessingTime     string
}

type ReportTrailer struct {
	DocumentCount int
	SuccessCount  int
	ErrorCount    int
}

type ReportDocument struct {
	ID           string    `xml:",attr"`
	FileName     string    `xml:",attr"`
	Status       string    // Succeeded, Failed
	ContentID    string    `xml:",omitempty"`
	ErrorCode    string    `xml:",omitempty"`
	ErrorMessage string    `xml:",omitempty"`
	Metadata     *Metadata `xml:",omitempty"`
}
