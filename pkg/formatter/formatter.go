package formatter

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"os"
	"strings"
)

type EncoderFunc func(data *string)

type Payload struct {
	encoders                            map[string]EncoderFunc
	encoderFlags                        []string
	Payload                             string
	defaultContentLength, ContentLength int
}

func NewPayload(payload string, flags []string) *Payload {
	encoders := map[string]EncoderFunc{
		"none": func(data *string) {
			return
		},
		"urlencode": func(data *string) {
			*data = url.QueryEscape(*data)
		},
		"base64": func(data *string) {
			*data = base64.StdEncoding.EncodeToString([]byte(*data))
		},
		"hex": func(data *string) {
			*data = hex.EncodeToString([]byte(*data))
		},
	}
	return &Payload{Payload: payload, encoders: encoders, encoderFlags: flags}
}

func CountDefaultContentLenght(requestfile string) int {
	request, _ := os.ReadFile(requestfile)
	var content string
	splited_request := strings.Split(strings.ReplaceAll(string(request), "\r\n", "\n"), "\n")
	for i := 0; i < len(splited_request); i++ {
		if len(splited_request[i]) == 0 {
			dataslice := splited_request[i:]
			content = strings.Join(strings.Fields(strings.Join(dataslice[:], "")), "")
			break
		}
	}
	return len(content) - len("{{.Payload}}")
}

func (p *Payload) AddDefaultContentLenght(lenght int) {
	p.defaultContentLength = lenght
}

func (p *Payload) CreatePayload() {
	for _, value := range p.encoderFlags {
		p.encoders[value](&p.Payload)
	}
	p.ContentLength = p.defaultContentLength + len(p.Payload)
}
