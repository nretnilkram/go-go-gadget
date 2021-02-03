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
        "github.com/nretnilkram/go-go-gadget/words"
)

func help() {
  fmt.Print(`
go-go-gadget is a set of command line tools.

Options:
  --help            Display go-go-gadget help
    -h, help
  --kubernetes      Take a string and k8s-ify it
    -k, --k8s, k8s
  --password        Return a password of desired length
    -p, --pw, pw, password
  --reverse         Take a string and reverse it
    -r, reverse
  --symsub          Take a string and substitute some chars to symbols
    -s, symsub
  --time            Display the current time
    -t, time
  --words           Return a desired number of words (english)
    -w, words

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
    options := []string{"help", "k8s", "password", "reverse", "symsub", "time", "words"}
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
    weight := pswd.PasswordWeight { Lower: 4, Upper: 3, Digit: 3, Symbol: 2, }
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
  case "symsub", "-s", "--symsub":
    if (len(params) > 0) {
      fmt.Println(strtwist.SymbolSubstitution(params[0]))
    } else {
      fmt.Print("Enter string to be subbed: ")
      str := prompt()
      fmt.Println(strtwist.SymbolSubstitution(str))
    }
  case "time", "-t", "--time":
    fmt.Println(timeStamp())
  case "words", "-w", "--words":
    weight := words.WordSetWeight { Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1, }
    if (len(params) > 0) {
      length, _ := strconv.Atoi(params[0])
      fmt.Println(words.Words(length, weight))
    } else {
      fmt.Print("How many words would you like: ")
      str := prompt()
      length, _ := strconv.Atoi(str)
      fmt.Println(words.Words(length, weight))
    }
  default:
    help()
  }
}
