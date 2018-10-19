package selpg

import (
	"os"
)

type Selpg struct {
	Begin int
	End int
	/* false for static line number */
	PageType bool
	Length int
	Destination string
	Src string
	data []string
	Logfile *os.File
}