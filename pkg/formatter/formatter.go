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
	DefaultContentLength, ContentLength int
}

func NewPayload(payload string, flags []string) *Payload {
	encoders := map[string]EncoderFunc{
		// "none" is needed in case the encoder is not passed, because then we get an error
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

// This function is needed because fasthttp does not calculate the Content-length header.
func CountDefaultContentLength(requestfile string) int {
	request, _ := os.ReadFile(requestfile)
	var content string
	splited_request := strings.Split(strings.ReplaceAll(string(request), "\r\n", "\n"), "\n")
	for i := 0; i < len(splited_request); i++ {
		if len(splited_request[i]) == 0 { // 0 - no data, just empty string
			dataslice := splited_request[i:]
			content = strings.Join(strings.Fields(strings.Join(dataslice[:], "")), "")
			break
		}
	}
	return len(content) - len("{{.Payload}}")
}

func (p *Payload) CreatePayload() {
	for _, value := range p.encoderFlags {
		p.encoders[value](&p.Payload)
	}
	p.ContentLength = p.DefaultContentLength + len(p.Payload)
}
