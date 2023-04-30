package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"strconv"
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

func GetResult() string {
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

	result := fmt.Sprintf("%s - %d - %s - %s", SelectedFormat, len(bytes), elapsedSerialization, elapsedDeserialization)
	log.Printf("result: %s", result)

	return result
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
	port, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("incorrect port is provided")
	}
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
	
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP:[]byte{0,0,0,0},Port:port,Zone:""})
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, _ := ServerConn.ReadFromUDP(buf)
		log.Printf("receive bytes: %s", string(buf[0:n]))
		if string(buf[0:n]) == "get_result" {
            result := GetResult()
			ServerConn.WriteTo([]byte(result), addr)
		}
	}
}
