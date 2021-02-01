package main

import (
        "bufio"
        "fmt"
        "os"
        "strings"
        "time"

        "github.com/nretnilkram/go-go-gadget/strtwist"
)

func help() {
  fmt.Print(`
go-go-gadget is a set of command line tools.

Options:
  help        Display go-go-gadget help
  k8s         Take a string and k8s-ify it
  reverse     Take a string and reverse it
  time        Display the current time

`)
}

func timeStamp() string {
  t := time.Now()
  return t.Format("06/01/02 15:04:05")
}

func getArgs() (string, []string) {
  executable := os.Args[0]
  parameters := os.Args[1:]
  return executable, parameters
}

func prompt() string {
  reader := bufio.NewReader(os.Stdin)
  text, _ := reader.ReadString('\n')
  // convert CRLF to LF
  text = strings.Replace(text, "\n", "", -1)
  return text
}

func main() {
  _, params := getArgs()

  var choice string

  if (len(params) > 0) {
    choice, params = params[0], params[1:]
  } else {
    options := []string{"help", "k8s", "reverse", "time"}
    fmt.Print("Options ", options, ": ")
    choice = prompt()
  }

  switch choice {
  case "help":
    help()
  case "k8s":
    if (len(params) > 0) {
      fmt.Println(strtwist.K8s(params[0]))
    } else {
      fmt.Print("Enter string to be K8s: ")
      str := prompt()
      fmt.Println(strtwist.K8s(str))
    }
  case "reverse":
    if (len(params) > 0) {
      fmt.Println(strtwist.Reverse(params[0]))
    } else {
      fmt.Print("Enter string to be reversed: ")
      str := prompt()
      fmt.Println(strtwist.Reverse(str))
    }
  case "time":
    fmt.Println(timeStamp())
  default:
    help()
  }
}
