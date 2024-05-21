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

type PkgHeader struct {
	SubmissionDate string
	SubmissionTime string
	Source         string
}

type PkgFooter struct {
	RequestCount int16
}

type Pkg struct {
	XMLName  xml.Name `xml:"Package"`
	ID       string   `xml:"ID,attr"`
	Header   PkgHeader
	Requests []Request `xml:"Requests>Request"`
	Footer   PkgFooter
}
