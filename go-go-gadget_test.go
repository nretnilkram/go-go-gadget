package main

import (
        "fmt"
        "os/exec"
        "strings"
        "testing"
)

func TestGoGoGadgetReverse(t *testing.T) {
  cases := []struct {
    in, want string
  }{
    {"Hello, world", "dlrow ,olleH"},
    {"Hello, 世界", "界世 ,olleH"},
    {"", ""},
    {"A", "A"},
    {"aBc", "cBa"},
    {"Mark", "kraM"},
  }
  for _, c := range cases {
    cmd := exec.Command("go-go-gadget", "reverse", c.in)
    out, err := cmd.CombinedOutput()
    got := strings.TrimSuffix(string(out), "\n") // because out is []byte
    if err != nil || got != c.want {
      fmt.Println(got, out, err)
      t.Errorf("go-go-gadget reverse %q == %q, want %q", c.in, got, c.want)
    }
  }
}

func TestGoGoGadgetK8s(t *testing.T) {
  cases := []struct {
    in, want string
  }{
    {"Kubernetes", "K8s"},
    {"kubernetes", "k8s"},
    {"kuberneteS", "k8S"},
    {"KUBERNETES", "K8S"},
    {"Mark", "M2k"},
    {"", ""},
    {"xy", "xy"},
    {"xyz", "x1z"},
    {"one two three", "o1e t1o t3e"},
  }
  for _, c := range cases {
    cmd := exec.Command("go-go-gadget", "k8s", c.in)
    out, err := cmd.CombinedOutput()
    got := strings.TrimSuffix(string(out), "\n") // because out is []byte with a new line
    if err != nil || got != c.want {
      fmt.Println(got, out, err)
      t.Errorf("go-go-gadget k8s %q == %q, want %q", c.in, got, c.want)
    }
  }
}
