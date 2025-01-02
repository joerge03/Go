package metadata

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
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
	
	fmt.Printf("token %v\n ", token)	
	v, ok := OfficeVersions[token[0]]
	if !ok {
		return "Unknown"
	}
	return v	
}


func NewProperties(r *zip.Reader) (*OfficeCoreProperty, *OfficeAppProperty, error){
	coreProp := new(OfficeCoreProperty)
	appProp := new(OfficeAppProperty)

	for _,f := range r.File {
		switch f.Name {			
		case "docProps/core.xml":
			if err := process(f,coreProp); err != nil {
				return nil,nil,err
			} 
			case "docProps/app.xml":
				if err := process(f,appProp); err != nil {
					return nil,nil,err
			}
		default:
			continue
		}
	}
	return coreProp,appProp,nil
}


func process(f *zip.File, prop any) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}

	defer rc.Close()
	
	err = xml.NewDecoder(rc).Decode(prop)
	if err != nil {
		return err
	}
	return nil
}