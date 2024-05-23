package formatter

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"
)

type EncoderFunc func(data string) string

type Payload struct {
	encoders map[string]EncoderFunc
	payload  string
}

func NewPayload(payload string) *Payload {
	encoders := map[string]EncoderFunc{
		"none": func(data string) string {
			return data
		},
		"urlencode": func(data string) string {
			return url.QueryEscape(data)
		},
		"base64": func(data string) string {
			return base64.StdEncoding.EncodeToString([]byte(data))
		},
		"hex": func(data string) string {
			return hex.EncodeToString([]byte(data))
		},
	}
	return &Payload{payload: payload, encoders: encoders}
}

func (p *Payload) CreatePayload(format string) string {
	splited_mask := strings.Split(strings.ReplaceAll(format, " ", ""), ",")
	for _, value := range splited_mask {
		p.payload = p.encoders[value](p.payload)
	}
	return p.payload
}
