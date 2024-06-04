# 🔥 WebFuzz - a very simple webfuzzer 
Written in go webfuzzer with multithreading support, fasthttp is used to handle http requests 

To work with webfuzz you need to upload a request to a file and use a template to specify where the payload from the dictionary will be substituted, it expands the possibility for phasing, from postdata to http-headers, also webfuzz itself supports various encoders that can help in testing. 

## Installation 
It is possible to download both ready binary file and build from source, but I recommend to download the binary file for your operating system, for example here is an example of downloading for Linux_x86_64
```sh 
wget https://github.com/LSD00/webfuzz/releases/download/v1.0.0/webfuzz_Linux_x86_64.tar.gz
```
## Usage 
To use webfuzz you need to write your raw http-requests in a txt file with a template, for example where None is the absence of any encoder : 
```
GET /{{ .Payload }} HTTP/1.1
Host: testphp.vulnweb.com
User-Agent: curl/7.81.0
Accept: */*


```
And the command to run would look like this, where r.txt is the query and w.txt is the dictionary with the directories or payload : 
```sh
webfuzz -d testphp.vulnweb.com -r r.txt -w w.txt --no-tls
```

There is also support for regular expressions, which can help in finding vulnerabilities such as sql injection, xss, os command injection, etc. 


It is also possible to filter inappropriate http status codes, which can help when testing for weak passwords (think Bruteforce) or when searching directories. 
