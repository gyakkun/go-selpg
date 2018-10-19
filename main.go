package main

import "flag"
import "os"
import "github.com/txzdream/serviceCourse/selpg/lib/selpg"
import "fmt"
import "time"

func main() {
	// ----定义flag
	flag.Usage = func() {
		fmt.Printf("Usage of seplg:\n")
		fmt.Printf("seplg -s num1 -e num2 [-f -l num3 -d str1 file]\n")
		flag.PrintDefaults()
	}
	start := flag.Int("s", -1, "Start of the page")
	end := flag.Int("e", -1, "End of the page")
	pagetype := flag.Bool("f", false, "If the page has static number of lines")
	length := flag.Int("l", -1, "the number of lines of every page")
	destination := flag.String("d", "", "the destination to write")
	flag.Parse()

	// ----处理输入错误
	/*开始结束页错误*/
	if *start <= 0 || *end <= 0 || *end < *start {
		fmt.Fprintf(os.Stderr, "Error: Invalid start, end page or line number. Use selpg -help to know more.\n")
		os.Exit(0)
	}
	/*同时存在-f和-l参数*/
	if *pagetype != false && *length != -1 {
		fmt.Fprintln(os.Stderr, "Error: Conflict flags: -f and -l")
		os.Exit(0)
	}
	/*设置行数默认值*/
	if *pagetype == false && *length <=0 {
		fmt.Println("Use 72 lines per page as default.")
		*length = 72
	}
	var src string
	if len(flag.Args()) == 1 {
		src = flag.Args()[0]
	} else if len(flag.Args()) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too much argument. Use selpg -help to know more.\n")
		os.Exit(0)
	} else {
		src = ""
	}

	// ----实例化对象
	data := selpg.Selpg{
		Begin: *start,
		End: *end,
		PageType: *pagetype,
		Length: *length,
		Destination: *destination,
		Src: src,
	}
	Logfile, err := os.OpenFile("log/log.txt", os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	Logfile.WriteString(time.Now().String())
	Logfile.WriteString("\n")

	// ----运行
	// 因为我不知道类似java的切片怎么去用，所以只能这种很丑的代码去完成log操作
	data.Read(Logfile)
	if data.Destination == "" {
		data.Write(Logfile)
	} else {
		data.Print(Logfile)
	}
	Logfile.Close()
}
