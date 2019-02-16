package main

import (
  "os"
  "os/exec"
  "flag"
  "fmt"
  "time"
  "log"
  "crypto/sha256"
  "io/ioutil"
  "strings"
)

const (
  // DefaultStoreDirectory is the default name of the store directory
  DefaultStoreDirectory = ".cacheexec"
)

// storePath tries to find a sensible place for the cache store
func storePath() string {
  home := os.Getenv("HOME")
  if home == "" {
    return DefaultStoreDirectory
  }
  return home + "/" + DefaultStoreDirectory
}

// readFromPipe executes a command and returns its output. If any
// errors are encountered, it exits with an error status.
func readFromPipe(args []string) ([]byte) {
  cmd := exec.Command(args[0], args[1:]...)
  pipe,err := cmd.StdoutPipe()
  if err != nil {
    log.Fatalf("error opening pipe from %q: %v", args[0], err)
  }
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
  data, err := ioutil.ReadAll(pipe)
  if err != nil {
    log.Fatalf("error reading from pipe from %q: %v", args[0], err)
  }
	if err := cmd.Wait(); err != nil {
    log.Fatalf("error waiting for %q: %v", args[0], err)
	}
  return data
}

// commandToKey turns a command+args into a store key
func commandToKey(args []string) string {
  digest := sha256.Sum256([]byte(strings.Join(args, " ")))
  return fmt.Sprintf("%x", digest)
}

func emit(data []byte) {
  n,err := os.Stdout.Write(data)
  if err != nil {
    log.Fatalf("could only write %d bytes to stdout. please help: %v", n, err)
    // CBF implementing a Write loop here
  }
  os.Exit(0)
}

func main() {
  ttl := flag.Duration("ttl", 8*time.Hour, "maximum time before cache invalidation")
  storePath := flag.String("storepath", storePath(), "location of cache store")
  flag.Parse()
  if flag.NArg() < 1 {
    log.Fatalf("A command must be specified")
  }

  args := flag.Args()
  store := NewStore(*storePath)
  key := commandToKey(args)
  if data,ok := store.Get(key,*ttl); ok {
    // Got valid cached data, just write it to stdout and exit
    emit(data)
  }
  // Nothing cached, or trouble reading from the cache. Refresh
  data := readFromPipe(args)
  store.Put(key, data)
  emit(data)
}
