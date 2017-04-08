package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"github.com/aphistic/docfs/dfs"
)

const (
	mountRoot = "/home/aphistic/tmp/docfs"
	docRoot   = "/home/aphistic/tmp/docroot"
)

func main() {
	fuse.Debug = func(msg interface{}) { fmt.Println(msg) }

	fs, err := dfs.NewDocFS(docRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "docfs root '%s' could not be opened: %s\n",
			docRoot, err)
		os.Exit(1)
	}
	defer fs.Close()

	c, err := fuse.Mount(mountRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't mount fs: %s\n", err)
		os.Exit(1)
	}
	defer c.Close()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT)

	serveChan := make(chan error, 1)
	go func() {
		serveChan <- fusefs.Serve(c, fs)
	}()

	select {
	case err := <-serveChan:
		if err != nil {
			fmt.Printf("Error running serve: %s\n", err)
		}
	case sig := <-sigChan:
		fmt.Printf("Signal %s received, stopping\n", sig)
		fuse.Unmount(mountRoot)
	}

	time.AfterFunc(2*time.Second, func() {
		fmt.Printf("Exiting timed out, exiting with error\n")
		os.Exit(1)
	})
}
