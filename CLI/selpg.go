package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	// "strings"

	"github.com/spf13/pflag"
)

type Selpg_Args struct {
	start_page *int
	end_page   *int
	// page_len    int
	page_len    *int
	page_type_f *bool
	output      *string
	input       string
	dest_print  *string
}

const INT_MAX = int(^uint(0) >> 1)

var selpg_args Selpg_Args
var print_pages int = 0

func initial() {
	// selpg_args.page_len = 72
	// 添加shorthand参数（去掉方法后面的字母P即取消shorthand参数）
	selpg_args.start_page = pflag.IntP("start_page", "s", -1, "start page")
	selpg_args.end_page = pflag.IntP("end_page", "e", -1, "end page")
	selpg_args.page_type_f = pflag.BoolP("form_feed", "f", false, "delimited by form feeds")
	selpg_args.page_len = pflag.IntP("limit", "l", 72, "delimited by fixed page length")
	selpg_args.output = pflag.StringP("output", "o", "", "output filename")
	selpg_args.dest_print = pflag.StringP("dest", "d", "", "target printer")
	// 另外一种写法
	// pflag.IntVarP(selpg_args.start_page, "start_page", "s", -1, "start page")
	// pflag.IntVarP(selpg_args.end_page, "end_page", "e", -1, "end page")
	// pflag.BoolVarP(selpg_args.page_type_f, "form_feed", "f", false, "delimited by form feeds")
	// pflag.IntVarP(selpg_args.page_len, "limit", "l", 72, "delimited by fixed page length")
	// pflag.StringVarP(selpg_args.output, "dest", "d", "", "output filename")

	selpg_args.input = ""
	pflag.Usage = usage
	pflag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, os.Args[0]+" -sstart_page_num -eend_page_num [ -f | -llines_per_page ] [ -doutput ] [ input ]\n")
	fmt.Fprintln(os.Stderr, "[OPTIONS]:")
	fmt.Fprintln(os.Stderr, "   filename string       input filename")
	pflag.PrintDefaults()
}

func checkErr(err error) {
	if err != nil {
		// panic(err)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func handle_args() {
	if len(os.Args) < 3 || *selpg_args.start_page == -1 || *selpg_args.end_page == -1 {
		fmt.Println("You are expected to input the start_page_num and end_page_num which will be greater than zero")
		pflag.Usage()
		os.Exit(1)
	}
	/* judge the start_page_num and end_page_num */
	if *selpg_args.start_page <= 0 {
		checkErr(errors.New("The start_page_num can't be less than or equal to zero"))
	}
	if *selpg_args.start_page > INT_MAX {
		checkErr(errors.New("The start_page_num can't be greater than INT_MAX: " + strconv.Itoa(INT_MAX)))
	}
	if *selpg_args.end_page <= 0 {
		checkErr(errors.New("The end_page_num can't be less than or equal to zero"))
	}
	if *selpg_args.end_page > INT_MAX {
		checkErr(errors.New("The end_page_num can't be greater than INT_MAX: " + strconv.Itoa(INT_MAX)))
	}
	if *selpg_args.start_page > *selpg_args.end_page {
		checkErr(errors.New("The end_page_num can't be greater than start_page_num"))
	}

	/* judge the page_len */
	if *selpg_args.page_len <= 0 {
		checkErr(errors.New("The -l limit can't be less than or equal to zero"))
	}

	for _, arg := range os.Args[1:] {
		/* judge the option -f */
		if matched, _ := regexp.MatchString(`^-f`, arg); matched {
			if len(arg) > 2 {
				checkErr(errors.New(os.Args[0] + ": option should be only \"-f\""))
			}
		}
		/* judge the dest printer */
		if matched, _ := regexp.MatchString(`^-d`, arg); *selpg_args.dest_print == "" && matched {
			checkErr(errors.New(os.Args[0] + ": option -d requires a printer destination\n"))
		}
		/* judge the output file */
		if matched, _ := regexp.MatchString(`^-o`, arg); *selpg_args.output == "" && matched {
			checkErr(errors.New(os.Args[0] + ": option -o requires a output filename\n"))
		}
		/* store the input filename */
		if arg[0] != '-' {
			selpg_args.input = arg
			break
		}
	}
}

func readFile(w *io.PipeWriter, InputFile *os.File) {
	inputReader := bufio.NewReader(InputFile)
	delimit := '\n'
	pages := 0
	lines := 0
	inputString := ""
	if *selpg_args.page_type_f {
		delimit = '\f'
	}
	for {
		inputSubString, readerError := inputReader.ReadString(byte(delimit))
		inputSubString = strings.Replace(inputSubString, "\f", "", -1)
		inputString += inputSubString
		if delimit == '\n' {
			lines++
			if lines == *selpg_args.page_len {
				pages++
				lines = 0
				if pages >= *selpg_args.start_page && pages < *selpg_args.end_page {
					print_pages++
					fmt.Fprint(w, inputString+"\f")
					inputString = ""
				}
			}
		} else {
			pages++
			if pages >= *selpg_args.start_page && pages < *selpg_args.end_page {
				print_pages++
				fmt.Fprint(w, inputString+"\f")
				inputString = ""
			}
		}
		if pages >= *selpg_args.end_page {
			print_pages++
			fmt.Fprint(w, inputString)
			break
		}
		if readerError == io.EOF {
			if inputString != "" && inputString != "\n" && inputString != "\f" {
				print_pages++
				fmt.Fprint(w, inputString)
			}
			break
		}
	}
	w.Close()
}

func handle_input() {
	var InputFile *os.File
	if selpg_args.input != "" {
		var err error
		InputFile, err = os.Open(selpg_args.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred on opening the inputfile: %s\nDoes the file exist?\n", selpg_args.input)
			os.Exit(1)
		}
	} else {
		InputFile = os.Stdin
	}
	defer InputFile.Close()

	var OutputFile *os.File
	var cmd *exec.Cmd
	if *selpg_args.output != "" {
		var err error
		OutputFile, err = os.Create(*selpg_args.output)
		checkErr(err)
	} else if *selpg_args.dest_print != "" {
		fmt.Printf("Printer: %s\n", *selpg_args.dest_print)
		cmd = exec.Command("lp", "-d", *selpg_args.dest_print)
		// cmd = exec.Command("lpr", "-P", *selpg_args.dest_print)
	} else {
		OutputFile = os.Stdout
	}
	defer OutputFile.Close()

	r, w := io.Pipe()
	go readFile(w, InputFile)
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	if cmd != nil {
		cmd.Stdin = strings.NewReader(buf.String())
		var outErr bytes.Buffer
		cmd.Stderr = &outErr
		// var out bytes.Buffer
		// cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			checkErr(errors.New(fmt.Sprint(err) + ": " + outErr.String()))
		}
		// fmt.Printf("in all caps: %q\n", out.String())
	}
	if OutputFile != nil {
		OutputFile.WriteString(buf.String())
	}
}

func main() {
	initial()
	handle_args()
	handle_input()
	fmt.Println("Print a total of " + strconv.Itoa(print_pages) + " pages!")
}
