package workers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"sync"
	"text/template"
	"webfuzz/pkg/formatter"
	"webfuzz/pkg/reqgen"
)

type Options struct {
	Concurrents                              int
	TlsEnabled                               bool
	Regex, Domain, RequestFile, WordlistFile string
	InvalidCode, Encoders                    []string
}

type Pool struct {
	wordlists      []string
	default_length int
	req            *template.Template
	options        Options
}

func NewPool(options Options) (*Pool, error) {
	var pool Pool
	pool.options = options
	file, err := os.Open(options.WordlistFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		pool.wordlists = append(pool.wordlists, reader.Text())
	}
	pool.req, err = template.ParseFiles(options.RequestFile)
	pool.default_length = formatter.CountDefaultContentLength(options.RequestFile)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

func (p *Pool) worker(wg *sync.WaitGroup, chunked_wordlist []string) {
	defer wg.Done()
	re := regexp.MustCompile(p.options.Regex)
	for _, value := range chunked_wordlist {
		var buf bytes.Buffer
		payload := formatter.NewPayload(value, p.options.Encoders)
		payload.DefaultContentLength = p.default_length
		payload.CreatePayload()
		p.req.Execute(&buf, payload)
		worker := reqgen.NewWorker(p.options.TlsEnabled, bufio.NewReader(&buf), p.options.Domain)
		resp, err := worker.MakeRequest()

		// Io.WriteString is necessary because it immediately writes to stdout
		// with fmt.Println there are problems and it may not output all
		if err != nil {
			io.WriteString(os.Stdout, err.Error())
			return
		} else if !slices.Contains(p.options.InvalidCode, fmt.Sprint(resp.Status)) && re.MatchString(resp.BodyData) {
			var formatted_code string
			if resp.Status >= 200 && resp.Status < 300 {
				formatted_code = fmt.Sprintf("\033[92m%d\033[0m", resp.Status)
			} else if resp.Status >= 300 && resp.Status < 400 {
				formatted_code = fmt.Sprintf("\033[96m%d\033[0m", resp.Status)
			} else if resp.Status >= 400 {
				formatted_code = fmt.Sprintf("\033[91m%d\033[0m", resp.Status)
			}
			io.WriteString(os.Stdout, fmt.Sprintf("%s -> %s (%d)\r\n", value, formatted_code, resp.Len))
		}
	}
}

func (p *Pool) Fuzz() {
	var wg sync.WaitGroup
	chunks := len(p.wordlists) / p.options.Concurrents
	for i := 0; i < p.options.Concurrents; i++ {
		var chunked_wordlist []string
		if i != p.options.Concurrents-1 {
			chunked_wordlist = p.wordlists[i*chunks : (i+1)*chunks]
		} else {
			chunked_wordlist = p.wordlists[i*chunks:]
		}
		wg.Add(1)
		go p.worker(&wg, chunked_wordlist)
	}
	wg.Wait()
}
