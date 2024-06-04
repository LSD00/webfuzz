package main

import (
	"fmt"
	"strings"
	"webfuzz/pkg/workers"

	"github.com/alecthomas/kingpin/v2"
)

var (
	invalid, encoders string
	options           workers.Options
)

func main() {
	kingpin.Flag("domain", "target domain").Required().Short('d').StringVar(&options.Domain)
	kingpin.Flag("request", "file with http request").Required().Short('r').StringVar(&options.RequestFile)
	kingpin.Flag("wordlist", "file with wordlist").Required().Short('w').StringVar(&options.WordlistFile)
	kingpin.Flag("regex", "regex for search").Default(".+").StringVar(&options.Regex)
	kingpin.Flag("threads", "setting threads").Default("25").IntVar(&options.Concurrents)
	kingpin.Flag("tls", "is target use tls").Default("false").BoolVar(&options.TlsEnabled)
	kingpin.Flag("bad-codes", "invalid codes for fuzzing ex. 404, 503, 400").Default("404").StringVar(&invalid)
	kingpin.Flag("encoders", "encode payloads: urlencode, base64, hex").Default("none").Short('e').StringVar(&encoders)
	kingpin.Parse()
	options.InvalidCode = strings.Split(strings.ReplaceAll(invalid, " ", ""), ",")
	options.Encoders = strings.Split(strings.ReplaceAll(encoders, " ", ""), ",")
	worker, err := workers.NewPool(options)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(`
	 _    _      _     ______             
	| |  | |    | |   |  ___|            
	| |  | | ___| |__ | |_ _   _ ________
	| |/\| |/ _ \ '_ \|  _| | | |_  /_  /
	\  /\  /  __/ |_) | | | |_| |/ / / / 
	 \/  \/ \___|_.__/\_|  \__,_/___/___|
										 
										 
										 `)
	worker.Fuzz()
}
