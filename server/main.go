package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"io"
	"github.com/go-yaml/yaml"
)

const Times = 1000

var SelectedFormat string
var SerializeFunc func(Object) []byte
var DeserializeFunc func([]byte)

type StringMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Object struct {
	Name         string            `json:"name" xml:"name" yaml:"name"`
	Precision    float64           `json:"precision" xml:"precision" yaml:"precision"`
	Order        int               `json:"order" xml:"order" yaml:"order"`
	Results      []int             `json:"results" xml:"results" yaml:"results"`
	Collocations StringMap         `json:"collocations" xml:"collocations" yaml:"collocations"`
}

func GetDefaultObject() Object {
	obj := Object{
		Name:         "Vladimir",
		Precision:    0.234,
		Order:        100,
		Results:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		Collocations: make(map[string]string),
	}
	obj.Collocations["a"] = "b"
	obj.Collocations["c"] = "d"
	obj.Collocations["Hello"] = "world"
	obj.Collocations["little"] = "biiiiiiiiiiiiiiiiiiiiiig"
	return obj
}

type Responce struct {
	Result string `json:"result"`
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	obj := GetDefaultObject()
	bytes := SerializeFunc(obj)

	startSerialization := time.Now()
	for i := 0; i < Times; i++ {
		SerializeFunc(obj)
	}
	elapsedSerialization := time.Since(startSerialization) / Times

	startDeserialization := time.Now()
	for i := 0; i < Times; i++ {
		DeserializeFunc(bytes)
	}
	elapsedDeserialization := time.Since(startDeserialization) / Times

	resp := Responce{
		Result: fmt.Sprintf("%s - %d - %s - %s", SelectedFormat, len(bytes), elapsedSerialization, elapsedDeserialization),
	}

	log.Printf("request result: %s", resp.Result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func SerializeJson(obj Object) []byte {
	bytes, _ := json.Marshal(obj)
	return bytes
}

func DeserializeJson(bytes []byte) {
	var obj Object
	json.Unmarshal(bytes, &obj)
}

func SerializeYaml(obj Object) []byte {
	bytes, _ := yaml.Marshal(obj)
	return bytes
}

func DeserializeYaml(bytes []byte) {
	var obj Object
	yaml.Unmarshal(bytes, &obj)
}

func SerializeXml(obj Object) []byte {
	bytes, _ := xml.Marshal(obj)
	return bytes
}

func DeserializeXml(bytes []byte) {
	var obj Object
	xml.Unmarshal(bytes, &obj)
}

func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func (m *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = StringMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func main() {
	args := os.Args
	if len(args) != 3 {
		log.Fatalf("incorrect num of args provided: 3 is required")
	}
	port := args[1]
	SelectedFormat = args[2]

	switch SelectedFormat {
	case "json":
		SerializeFunc = SerializeJson
		DeserializeFunc = DeserializeJson
	case "xml":
		SerializeFunc = SerializeXml
		DeserializeFunc = DeserializeXml
	case "yaml":
		SerializeFunc = SerializeYaml
		DeserializeFunc = DeserializeYaml

	default:
		log.Fatalf("format %s is not supported", SelectedFormat)
	}
	http.HandleFunc("/get_result", GetResult)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("server stopped with error: %s", err)
	}
}
