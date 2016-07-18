package term;

import (
  "bufio"
  "os"
  "fmt"
	"github.com/mattn/go-colorable"
  "encoding/base64"
  "net/http"
  "io/ioutil"
  "sync"
  "time"
  "math/rand"
)

const (
  Reset = "\x1B[0m"
  Bold = "\x1B[1m"
  Dim = "\x1B[2m"
  Under = "\x1B[4m"
  Reverse = "\x1B[7m"
  Hide = "\x1B[8m"
  Clearscreen = "\x1B[2J"
  Clearline = "\x1B[2K"
  Black = "\x1B[30m"
  Red = "\x1B[31m"
  Green = "\x1B[32m"
  Yellow = "\x1B[33m"
  Blue = "\x1B[34m"
  Magenta = "\x1B[35m"
  Cyan = "\x1B[36m"
  White = "\x1B[37m"
  Bblack = "\x1B[40m"
  Bred = "\x1B[41m"
  Bgreen = "\x1B[42m"
  Byellow = "\x1B[43m"
  Bblue = "\x1B[44m"
  Bmagenta = "\x1B[45m"
  Bcyan = "\x1B[46m"
  Bwhite = "\x1B[47m"
  Newline = "\r\n\x1B[0m"
  ImageBASE64 = "\033]1337;File=name=%s;inline=1;width=100px;height=auto:%s\a\n"
  ImageURL = "\033]1338;url=%s;alt=%s"
  CURSORBACK = "\033[1D"
  CURSORHIDE = "\033[?25l"
  CURSORSHOW = "\033[?25h"
)

var ERR_MUTEX sync.Mutex
var MESSAGE_MUTEX sync.Mutex

var THROBBERS = [][]string{
      []string{"⠋","⠙","⠹","⠸","⠼","⠴","⠦","⠧","⠇","⠏"},
      []string{"⠋","⠙","⠚","⠞","⠖","⠦","⠴","⠲","⠳","⠓"},
      []string{"⠄","⠆","⠇","⠋","⠙","⠸","⠰","⠠","⠰","⠸","⠙","⠋","⠇","⠆"},
      []string{"⠋","⠙","⠚","⠒","⠂","⠂","⠒","⠲","⠴","⠦","⠖","⠒","⠐","⠐","⠒","⠓","⠋"},
      []string{"⠁","⠉","⠙","⠚","⠒","⠂","⠂","⠒","⠲","⠴","⠤","⠄","⠄","⠤","⠴","⠲","⠒","⠂","⠂","⠒","⠚","⠙","⠉","⠁"},
      []string{"⠈","⠉","⠋","⠓","⠒","⠐","⠐","⠒","⠖","⠦","⠤","⠠","⠠","⠤","⠦","⠖","⠒","⠐","⠐","⠒","⠓","⠋","⠉","⠈"},
      []string{"⠁","⠁","⠉","⠙","⠚","⠒","⠂","⠂","⠒","⠲","⠴","⠤","⠄","⠄","⠤","⠠","⠠","⠤","⠦","⠖","⠒","⠐","⠐","⠒","⠓","⠋","⠉","⠈","⠈","⠉"},
    }

var STD_OUT = bufio.NewWriter(colorable.NewColorableStdout())
var STD_ERR = bufio.NewWriter(colorable.NewColorableStderr())
var STD_IN = bufio.NewReader(os.Stdin)


func OutputError(message string) {
  if(message != "") {
    ERR_MUTEX.Lock()
    STD_ERR.WriteString(Bold + Red + message + Reset + "\n")
  	STD_ERR.Flush()
    ERR_MUTEX.Unlock()
  }
}

func Throbber() {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  throb := THROBBERS[r.Intn(len(THROBBERS))]

  i := r.Intn(len(throb))

  OutputMessage(" ")
  OutputMessage(CURSORHIDE)
  for {
    i = (i + 1) % len(throb)
    OutputMessage(CURSORBACK+throb[i])
    time.Sleep(100 * time.Millisecond)
  }
  OutputMessage(CURSORSHOW)
}

func OutputMessage(message string) {
  MESSAGE_MUTEX.Lock()
	STD_OUT.WriteString(message)
	STD_OUT.Flush()
  MESSAGE_MUTEX.Unlock()
}

func OutputImageUrl(url string,alt string) {
  client := http.Client{}
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    OutputError(fmt.Sprintf("Error - %s",err.Error()))
  }
  resp, doError := client.Do(req)
  defer resp.Body.Close()
  if doError != nil {
    OutputError("Error - "+doError.Error())
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    OutputError(fmt.Sprintf("Error - %s",err.Error()))
  }
  encodedData := base64.StdEncoding.EncodeToString([]byte(body))
  encodedName := base64.StdEncoding.EncodeToString([]byte(alt))
  OutputMessage(fmt.Sprintf(ImageBASE64,encodedName,encodedData))
}
