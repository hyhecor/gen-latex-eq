package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	version       = "0.0.0"
	build         = "0"
	versionString = func() string {
		return fmt.Sprintf("%s@%s", version, build)
	}
)

// Args is
type Args struct {
	help         bool
	version      bool
	convertorURL string
}

var (
	args Args
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s %s:\n", os.Args[0], versionString())
		flag.PrintDefaults()
	}

	flag.BoolVar(&args.help, "version", false, "version")
	flag.BoolVar(&args.help, "h", false, "help")
	flag.StringVar(&args.convertorURL, "convertor", "https://latex.codecogs.com/svg.latex?", "")

	flag.Parse() // flag.Parse

	if args.help {
		flag.Usage()
		os.Exit(1)
	}
	if args.version {
		flag.Usage()
		os.Exit(1)
	}

	for _, f := range flag.Args() {
		log.Println(f)
	}
}

func main() {

	sliceEq := make([]eq, 0, 100)
	var elementEq eq
	pos := 0
	scanner := bufio.NewScanner(os.Stdin)
LOOP_SCANNER:
	for scanner.Scan() {
		pos++
		s := scanner.Text()
		log.Printf("%04d: %s\n", pos, s)

		// ss := strings.SplitN(s, "=", 1)
		ss := simpleSplit(s, "=", 1)
		switch len(ss) {
		case 2:
			k := strings.TrimSpace(ss[0])
			v := strings.TrimSpace(ss[1])

			if ';' != []byte(k)[0] {
				elementEq = eq{}
				elementEq.Name = k
				elementEq.Equation = v
				sliceEq = append(sliceEq, elementEq)
			}
		default:
			continue LOOP_SCANNER
		}

	}
	log.Printf("(%-8d)==================================================\n", len(sliceEq))

	in := make(chan eq)
	out := make(chan error)

	// error handler
	go func() {
		defer close(out)

		select {
		case err := <-out:
			if nil != err {
				log.Panic(err)
			}
		}
	}()

	//feed
	go func() {
		defer close(in)
		for _, element := range sliceEq {
			in <- element
		}
	}()

	fnStartWorker := WorkerFactory(args.convertorURL, in, out)
	var wg sync.WaitGroup

	for index := 0; index < runtime.NumCPU(); index++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			defer log.Println("Close Worker", n)
			fnStartWorker(n)
		}(index)
	}

	wg.Wait()

}

func simpleSplit(s, sep string, n int) []string {

	out := make([]string, 0)
LOOP_FOR:
	for true {
		n--
		if n < 0 {
			break LOOP_FOR
		}
		idx := strings.Index(s, sep)
		if idx == -1 {
			return out
		}

		out = append(out, strings.TrimSpace(s[:idx]))
		s = s[idx+1:]
	}
	out = append(out, strings.TrimSpace(s))
	return out
}

type eq struct {
	Name     string
	Equation string
}

func (s eq) EquationURL() string {
	return url.PathEscape(s.Equation)
}

// WorkerFactory is
func WorkerFactory(url string, in <-chan eq, out chan<- error) (closure func(int)) {

	closure = func(id int) {

	FOR_LOOP:
		for eq := range in {
			log.Printf("w_id=%v name=%v eq=%v\n", id, eq.Name, eq.Equation)

			filename := eq.Name
			equationURL := url + eq.EquationURL()

			// 경로 만들기
			dir := filepath.Dir(filename)
			err := os.MkdirAll(dir, 0755)
			if nil != err {
				continue FOR_LOOP
			}

			// 파일 열기
			f, err := os.OpenFile(filename, os.O_CREATE, 0755)
			if nil != err {
				out <- errors.WithMessage(err, fmt.Sprintf("w_id=%v name=%v eq=%v url=%v", id, filename, eq.Equation, equationURL))
				continue FOR_LOOP
			}
			defer f.Close()

			// 요청
			response, err := http.Get(equationURL)
			if nil != err {
				out <- err
				continue FOR_LOOP
			}
			defer response.Body.Close()

			// 응답 파일에 저장
			_, err = io.Copy(f, response.Body)
			if nil != err {
				out <- err
				continue FOR_LOOP
			}
		}
	}
	return
}
