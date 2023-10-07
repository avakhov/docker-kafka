package main

import (
	"fmt"
	"time"
  "reflect"
	"github.com/samuel/go-zookeeper/zk"
  "os"
  "strings"
  "runtime"
  "errors"
)

func wrapError(in interface{}, a ...any) error {
	_, full, line, _ := runtime.Caller(1)
	root, err := os.Getwd()
	if err != nil {
		return errors.New("sys error:" + err.Error())
	}
	file := strings.Replace(full, root+"/", "", 1)
	prefix := fmt.Sprintf("%s:%d →\n", file, line)
	switch v := in.(type) {
	case error:
		return errors.New(prefix + v.Error())
	case string:
		return errors.New(fmt.Sprintf(prefix+v, a...))
	case nil:
		return nil
	default:
		return errors.New(prefix + " unknown in type")
	}
}

func doMain(zoo_url string, brokers []string) error {
	conn, _, err := zk.Connect([]string{zoo_url}, time.Second)
	if err != nil {
    return wrapError(err)
	}
	defer conn.Close()

  for {
    current, _, err := conn.Children("/brokers/ids")
    if err != nil {
      return wrapError(err)
    }
    fmt.Printf("current brokers: '%s'\n", strings.Join(current, ", "))
    if reflect.DeepEqual(current, brokers) {
      break
    }
    time.Sleep(1 * time.Second)
  }
  return nil
}

func main() {
  fmt.Printf("Starting kafka-wait-brokers\n")
  if len(os.Args) < 3 {
    fmt.Printf("Usage: ./kafka-wait-brokers <zookeeper_url> [kafka_broker1 kafka_broker2 ...]\n")
    os.Exit(1)
  }
  err := doMain(os.Args[1], os.Args[2:])
  if err != nil {
    fmt.Printf("ERROR: %s\n", err.Error())
    os.Exit(1)
  }
}
