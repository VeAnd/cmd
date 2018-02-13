package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sync"

	"github.com/VeAnd/pkg/ansicolor"
)

const Version = "0.0.1"

var (
	versionFlag  = flag.Bool("version", false, "show program's version number and exit")
	helpFlag     = flag.Bool("help", false, "show this help message and exit")
	blueFlag     = flag.String("blue", "", "display BLUE pattern in blue")
	redFlag      = flag.String("red", "", "display RED pattern in red")
	greenFlag    = flag.String("green", "", "display GREEN pattern in green")
	yellowFlag   = flag.String("yellow", "", "display YELLOW pattern in green")
	magentaFlag  = flag.String("magenta", "", "display MAGENTA pattern in magenta")
	cyanFlag     = flag.String("cyan", "", "display CYAN pattern in cyan")
	bblueFlag    = flag.String("bblue", "", "display BBLUE pattern in blue background")
	bredFlag     = flag.String("bred", "", "display BRED pattern in red background")
	bgreenFlag   = flag.String("bgreen", "", " display BGREEN pattern in green background")
	byellowFlag  = flag.String("byellow", "", " display BYELLOW pattern in yellow background")
	bmagentaFlag = flag.String("bmagenta", "", " display BMAGENTA pattern in magenta background")
	bcyanFlag    = flag.String("bcyan", "", " display BCYAN pattern in cyan background")
	templateFlag = flag.String("template", "", "color text on template pattern")
)

func init() {
	flag.BoolVar(versionFlag, "v", false, "show program's version number and exit")
	flag.BoolVar(helpFlag, "h", false, "show this help message and exit")
	flag.StringVar(blueFlag, "b", "", "display BLUE pattern in blue")
	flag.StringVar(redFlag, "r", "", "display RED pattern in red")
	flag.StringVar(greenFlag, "g", "", "display GREEN pattern in green")
	flag.StringVar(yellowFlag, "y", "", "display YELLOW pattern in yellow")
	flag.StringVar(magentaFlag, "m", "", "display MAGENTA pattern in magenta")
	flag.StringVar(cyanFlag, "c", "", "display CYAN pattern in cyan")
	flag.StringVar(bblueFlag, "B", "", "display BBLUE pattern in blue background")
	flag.StringVar(bredFlag, "R", "", "display BRED pattern in red background")
	flag.StringVar(bgreenFlag, "G", "", "display BGREEN pattern in green background")
	flag.StringVar(byellowFlag, "Y", "", "display BYELLOW pattern in yellow background")
	flag.StringVar(bmagentaFlag, "M", "", "display BMAGENTA pattern in magenta background")
	flag.StringVar(bcyanFlag, "C", "", "display BBCYAN pattern in cyan background")
	flag.StringVar(templateFlag, "t", "", "color text on template pattern")
}

type patternFunc struct {
	Pattern string
	FName   ansicolor.ColorFunc
}

func main() {
	/*
	   	--version                         show program's version number and exit
	   --help, -h                        show this help message and exit
	   --blue=BLUE, -b BLUE              display BLUE pattern in blue
	   --red=RED, -r RED                 display RED pattern in red
	   --green=GREEN, -g GREEN           display GREEN pattern in green
	   --yellow=YELLOW, -y YELLOW        display YELLOW pattern in yellow
	   --magenta=MAGENTA, -m MAGENTA     display MAGENTA pattern in magenta
	   --cyan=CYAN, -c CYAN              display CYAN pattern in cyan
	   --bblue=BBLUE, -B BBLUE           display BBLUE pattern in blue background
	   --bred=BRED, -R BRED              display BRED pattern in red background
	   --bgreen=BGREEN, -G BGREEN        display BGREEN pattern in green background
	   --byellow=BYELLOW, -Y BYELLOW     display BYELLOW pattern in yellow background
	   --bmagenta=BMAGENTA, -M BMAGENTA  display BMAGENTA pattern in magenta background
	   --bcyan=BCYAN, -C BCYAN           display BCYAN pattern in cyan background
	   --blink=BLINK, -K BLINK           display BLINK pattern blinking (not widely supported)
	   --bold=BOLD, -D BOLD              display BOLD pattern in bold
	   --undrln=UNDRLN, -u UNDRLN        display UNDRLN pattern underlined
	   --bisounours, -N                  display with random colors
	*/

	if len(os.Args) == 1 {
		fmt.Println("-version\t\t\tshow program's version number and exit")
		fmt.Println("-help,\t\t-h\t\tshow this help message and exit")
		fmt.Println("-blue=BLUE,\t-b BLUE\t\tdisplay BLUE pattern in blue")
		fmt.Println("-red=RED,\t-r RED\t	display RED pattern in red")
		fmt.Println("-green=GREEN,\t-g GREEN\tdisplay GREEN pattern in green")

		return
	}

	flag.Parse()

	if *versionFlag {
		fmt.Println("go-rainbow: ", Version)
		os.Exit(0)
	}

	var apf []patternFunc

	if blueFlag != nil {
		apf = append(apf, patternFunc{*blueFlag, ansicolor.BoldBlue})
	}
	if redFlag != nil {
		apf = append(apf, patternFunc{*redFlag, ansicolor.BoldRed})
	}
	if greenFlag != nil {
		apf = append(apf, patternFunc{*greenFlag, ansicolor.BoldGreen})
	}
	if yellowFlag != nil {
		apf = append(apf, patternFunc{*yellowFlag, ansicolor.BoldYellow})
	}
	if magentaFlag != nil {
		apf = append(apf, patternFunc{*magentaFlag, ansicolor.BoldMagenta})
	}
	if cyanFlag != nil {
		apf = append(apf, patternFunc{*cyanFlag, ansicolor.BoldCyan})
	}
	if bblueFlag != nil {
		apf = append(apf, patternFunc{*bblueFlag, ansicolor.BBlue})
	}
	if bredFlag != nil {
		apf = append(apf, patternFunc{*bredFlag, ansicolor.BRed})
	}
	if bgreenFlag != nil {
		apf = append(apf, patternFunc{*bgreenFlag, ansicolor.BGreen})
	}
	if byellowFlag != nil {
		apf = append(apf, patternFunc{*byellowFlag, ansicolor.BYellow})
	}
	if bmagentaFlag != nil {
		apf = append(apf, patternFunc{*bmagentaFlag, ansicolor.BMagenta})
	}
	if bcyanFlag != nil {
		apf = append(apf, patternFunc{*bcyanFlag, ansicolor.BCyan})
	}
	
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		colorStream(reader, apf)
	}

	

}

var mu = &sync.Mutex{}

func colorStream(r *bufio.Reader, apf []patternFunc) {
	for {
		line, err := r.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}

		var inStr = line
		for _, pfs := range apf {
			re := regexp.MustCompile(pfs.Pattern)
			strings := re.FindAllString(line, -1)
			var wg sync.WaitGroup
			for _, str := range strings {
				wg.Add(1)
				go func() {
					mu.Lock()
					inStr = re.ReplaceAllString(inStr, pfs.FName(str))
					mu.Unlock()
					wg.Done()
				}()
			}
			wg.Wait()

		}
		fmt.Print(inStr)
	}
}
