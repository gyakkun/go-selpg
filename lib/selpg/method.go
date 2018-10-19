package selpg

import (
	"fmt"
	"io"
	"bufio"
	"os"
)

// ----从stdin或文件读取输入并存储合适的范围
func (selpg *Selpg) Read(Logfile *os.File) {
	if selpg == nil {
		fmt.Fprintf(os.Stderr, "Error: Unknown error.\n")
		Logfile.WriteString("[error] Use null object\n")
		os.Exit(0)
	}
	var in io.Reader

	// ----确定内容来源
	if selpg.Src == "" {
		in = os.Stdin
	} else {
		var err error
		in, err = os.Open(selpg.Src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: No such file found. Please pass right path.\n")
			Logfile.WriteString("[error] Unknown file to be read\n")
			os.Exit(0)
		}
	}
	
	// ----读取内容
	scanner := bufio.NewScanner(in)
	if selpg.PageType == false {
		cnt := 0
		for scanner.Scan() {
			line := scanner.Text()
			if cnt / selpg.Length + 1 >= selpg.Begin && cnt / selpg.Length + 1 <= selpg.End {
				selpg.data = append(selpg.data, line)
			}
			cnt++
		}
	} else {
		cnt := 1
		onComma := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			for i := 0; i < len(data); i++ {
				if data[i] == '\f' {
					return i + 1, data[:i], nil
				}
			}
			if atEOF {
				return 0, data, bufio.ErrFinalToken
			} else {
				return 0, nil, nil
			}
		}
		scanner.Split(onComma)
		for scanner.Scan() {
			line := scanner.Text()
			if cnt >= selpg.Begin && cnt <= selpg.End {
				selpg.data = append(selpg.data, line)
			}
			cnt++
		}
	}
	Logfile.WriteString("[info]  Read data finished\n")
}

// ----输出内容到stdout
func (selpg *Selpg) Write(Logfile *os.File) {
	if selpg == nil {
		fmt.Fprintf(os.Stderr, "Error: Unknown error.\n")
		Logfile.WriteString("[error] Use null object\n")
		os.Exit(0)
	}
	for i := 0; i < len(selpg.data); i++ {
		fmt.Fprintln(os.Stdout, selpg.data[i])
	}
	Logfile.WriteString("[info]  Write data finished\n")
}

// ----连接到打印机
func (selpg *Selpg) Print(Logfile *os.File) {
	if selpg.Destination != "" {
		file, err := os.Create(selpg.Destination)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: Can not create such file")
			Logfile.WriteString("[error] can not create file as -d argument\n")
			os.Exit(0)
		}
		for i := 0; i < len(selpg.data); i++ {
			file.WriteString(selpg.data[i])
			file.WriteString("\n")
		}
	}
	Logfile.WriteString("[info]  Print data finished\n")
}