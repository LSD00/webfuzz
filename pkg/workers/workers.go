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
	TlsEnabled  bool
	Regex       string
	InvalidCode []string
}

type Pool struct {
	concurrents      int
	wordlists, flags []string
	deflenght        int
	req              *template.Template
	Options          Options
}

func NewPool(wordlist, reqfile string, concurrents int) (*Pool, error) {
	var pool Pool
	pool.concurrents = concurrents
	file, err := os.Open(wordlist)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		pool.wordlists = append(pool.wordlists, reader.Text())
	}
	pool.req, err = template.ParseFiles(reqfile)
	pool.deflenght = formatter.CountDefaultContentLenght(reqfile)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

func (p *Pool) AddFlags(flags []string) {
	p.flags = flags
}

func (p *Pool) worker(wg *sync.WaitGroup, wordlist []string, domain string) {
	defer wg.Done()
	re := regexp.MustCompile(p.Options.Regex)
	for _, value := range wordlist {
		var buf bytes.Buffer
		payload := formatter.NewPayload(value, p.flags)
		payload.AddDefaultContentLenght(p.deflenght)
		payload.CreatePayload()
		p.req.Execute(&buf, payload)
		worker := reqgen.NewWorker(p.Options.TlsEnabled, bufio.NewReader(&buf), domain)
		resp, err := worker.MakeRequest()
		io.WriteString(os.Stdout, buf.String())
		if err != nil {
			io.WriteString(os.Stdout, err.Error())
			return
		} else if !slices.Contains(p.Options.InvalidCode, fmt.Sprint(resp.Status)) && re.MatchString(resp.BodyData) {
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

func (p *Pool) Fuzz(domain string) {
	var wg sync.WaitGroup
	chunks := len(p.wordlists) / p.concurrents
	for i := 0; i < p.concurrents; i++ {
		var chunked_wordlist []string
		if i != p.concurrents-1 {
			chunked_wordlist = p.wordlists[i*chunks : (i+1)*chunks]
		} else {
			chunked_wordlist = p.wordlists[i*chunks:]
		}
		wg.Add(1)
		go p.worker(&wg, chunked_wordlist, domain)
	}
	wg.Wait()
}
