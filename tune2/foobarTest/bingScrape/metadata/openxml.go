package metadata

import (
	"encoding/xml"
	"strings"
)


type OfficeCoreProperty struct{
	XMLName xml.Name `xml:"coreProperties"`
	Creator string `xml:"creator"`
	LastModifiedBy string `xml:"lastModifiedBy"`
}

type OfficeAppProperty struct{
	XMLName xml.Name `xml:"Properties"`
	Application string `xml:"Application"`	
	Company string `xml:"Company"`	
	Version string `xml:"AppVersion"`	
}

var OfficeVersions = map[string]string{
	"16": "2016",
	"15": "2013",
	"14": "2010",
	"12": "2007",
	"11": "2003",
}


func (a *OfficeAppProperty) GetMajorVersion() string{
	token := strings.Split(a.Version, ".")

	if len(token) < 2 {
		return "Unknown"
	}
	
	// v, ok := OfficeVersions[token[0]]
	return ""
}