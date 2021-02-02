package main

import (
        "bufio"
        "fmt"
        "os"
        "strings"
        "strconv"
        "time"

        "github.com/nretnilkram/go-go-gadget/pswd"
        "github.com/nretnilkram/go-go-gadget/strtwist"
)

func help() {
  fmt.Print(`
go-go-gadget is a set of command line tools.

Options:
  --help            Display go-go-gadget help
    -h, help
  --kubernetes      Take a string and k8s-ify it
    -k, --k8s, k8s
  --password         Take a string and reverse it
    -p, --pw, pw, password
  --reverse         Take a string and reverse it
    -r, reverse
  --time            Display the current time
    -t, time

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
    options := []string{"help", "k8s", "password", "reverse", "time"}
    fmt.Print("Options ", options, ": ")
    choice = prompt()
  }

  switch choice {
  case "help", "--help", "-h":
    help()
  case "k8s", "-k", "--kkubernetes", "--k8s":
    if (len(params) > 0) {
      fmt.Println(strtwist.K8s(params[0]))
    } else {
      fmt.Print("Enter string to be K8s: ")
      str := prompt()
      fmt.Println(strtwist.K8s(str))
    }
  case "password", "pw", "--pw", "-p", "--password":
    weight := PasswordWeight { lower: 4, upper: 3, digit: 3, symbol: 2, }
    if (len(params) > 0) {
      length, _ := strconv.Atoi(params[0])
      fmt.Println(pswd.Password(length, weight))
    } else {
      fmt.Print("How long should the password be: ")
      str := prompt()
      length, _ := strconv.Atoi(str)
      fmt.Println(pswd.Password(length, weight))
    }
  case "reverse", "-r", "--reverse":
    if (len(params) > 0) {
      fmt.Println(strtwist.Reverse(params[0]))
    } else {
      fmt.Print("Enter string to be reversed: ")
      str := prompt()
      fmt.Println(strtwist.Reverse(str))
    }
  case "time", "-t", "--time":
    fmt.Println(timeStamp())
  default:
    help()
  }
}
