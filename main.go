package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"crypto/tls"
)

var (
  Info = Teal
  Warn = Yellow
  Fata = Red
  Contype = Green
)

var (
  Black   = Color("\033[1;30m%s\033[0m")
  Red     = Color("\033[1;31m%s\033[0m")
  Green   = Color("\033[1;32m%s\033[0m")
  Yellow  = Color("\033[1;33m%s\033[0m")
  Purple  = Color("\033[1;34m%s\033[0m")
  Magenta = Color("\033[1;35m%s\033[0m")
  Teal    = Color("\033[1;36m%s\033[0m")
  White   = Color("\033[1;37m%s\033[0m")
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	flag.Parse()

	var input io.Reader
	input = strings.NewReader(strings.Join(flag.Args(), "\n"))
	if flag.NArg() == 0 {
		input = os.Stdin
	}

	sc := bufio.NewScanner(input)

	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		var payloadArr [3]string
		payloadArr[0] = "/(A(%22testabcd))"
		payloadArr[1] = "/(A(%3Ctestabcd))"
		payloadArr[2] = "/(A(%27testabcd))"
		for i := 0; i < len(payloadArr); i++ {
			finalurl := ""
			if u.Path == "" {
			finalurl = u.Scheme + "://" + u.Host + payloadArr[i] + "/"
			} else {
			finalurl = u.Scheme + "://" + u.Host + payloadArr[i] + u.Path
			}
			resp, err := http.Get(finalurl)
			if err != nil {
				if !strings.Contains(err.Error(), "no such host") {
					fmt.Fprintln(os.Stderr, err)
					continue
				} 
				continue
			}

			defer resp.Body.Close()

			serverheader := ""
			contenttype := ""
			for k, v := range resp.Header {
				if strings.ToLower(k) == "server" {
					serverheader = strings.ToLower(v[0])
				} else if strings.ToLower(k) == "content-type" {
					contenttype = strings.ToLower(v[0])
				}
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			body := string(b)

			re, err := regexp.Compile(`.*\"testabcd\)|.*\'testabcd\)|.*\<testabcd\)`) //"
			if err != nil {
				fmt.Fprintf(os.Stderr, "regexp compile error: %s", err)
			}

			matches := re.FindAllStringSubmatch(body, -1)

			for _, m := range matches {
				matched, _ := regexp.MatchString(`\\"testabcd`, m[0]) //"
				if !matched {
				
				//Colors for Quality of life
					if i == 0 {
					fmt.Printf(Warn("SUCCESS "))
						if contenttype == "application/xml" {
							fmt.Printf(Contype("(XML) "))
						} else if contenttype == "application/json" {
							fmt.Printf(Contype("(JSON) "))
						}
					} else if i == 1 {
					fmt.Printf(Fata("SUCCESS "))
					} else if i == 2 {
					fmt.Printf(Info("SUCCESS "))
						if contenttype == "application/xml" {
							fmt.Printf(Contype("(XML) "))
						} else if contenttype == "application/json" {
							fmt.Printf(Contype("(JSON) "))
						}
					}
					fmt.Printf("reflected for %s on %s for value %s\n", finalurl, serverheader, m[0])
				}
			}
		}
	}
}

func Color(colorString string) func(...interface{}) string {
  sprint := func(args ...interface{}) string {
    return fmt.Sprintf(colorString,
      fmt.Sprint(args...))
  }
  return sprint
}