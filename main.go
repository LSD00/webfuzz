package main

import (
	"fmt"
	"strings"
	"webfuzz/pkg/workers"

	"github.com/alecthomas/kingpin/v2"
)

var (
	domain, reqfile, wordlist, invalid, flags string
	th                                        int
	opt                                       workers.Options
)

func main() {
	kingpin.Flag("domain", "target domain").Required().Short('d').StringVar(&domain)
	kingpin.Flag("request", "file with http request").Required().Short('r').StringVar(&reqfile)
	kingpin.Flag("wordlist", "file with wordlist").Required().Short('w').StringVar(&wordlist)
	kingpin.Flag("regex", "regex for search").Default(".+").StringVar(&opt.Regex)
	kingpin.Flag("threads", "setting threads").Default("25").IntVar(&th)
	kingpin.Flag("tls", "is target use tls").Default("false").BoolVar(&opt.TlsEnabled)
	kingpin.Flag("bad-codes", "invalid codes for fuzzing ex. 404, 503, 400").Default("404").StringVar(&invalid)
	kingpin.Flag("encoders", "encode payloads: urlencode, base64, hex").Default("none").Short('e').StringVar(&flags)
	kingpin.Parse()
	opt.InvalidCode = strings.Split(strings.ReplaceAll(invalid, " ", ""), ",")
	worker, err := workers.NewPool(wordlist, reqfile, th)
	if err != nil {
		fmt.Println(err)
	}
	worker.AddFlags(strings.Split(strings.ReplaceAll(flags, " ", ""), ","))
	worker.Options = opt
	fmt.Println(`
	 _    _      _     ______             
	| |  | |    | |   |  ___|            
	| |  | | ___| |__ | |_ _   _ ________
	| |/\| |/ _ \ '_ \|  _| | | |_  /_  /
	\  /\  /  __/ |_) | | | |_| |/ / / / 
	 \/  \/ \___|_.__/\_|  \__,_/___/___|
										 
										 
										 `)
	worker.Fuzz(domain)
}
