package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"gopkg.in/fsnotify.v1"
)

// some systems say `foo' and some say 'foo'
var prereqDebugLine = regexp.MustCompile(`Prerequisite .([^']+)' is (older|newer) than target .([^']+)'.`)

var watchFailedColor = color.New(color.FgMagenta, color.Bold)
var changeColor = color.New(color.FgYellow, color.Bold)
var makingColor = color.New(color.FgBlue, color.Bold)
var makeSucceededColor = color.New(color.FgGreen, color.Bold)
var makeFailedColor = color.New(color.FgRed, color.Bold)

type Harry struct {
	MakeArgs []string

	watcher *fsnotify.Watcher
}

func newHarry(makeArgs []string) *Harry {
	return &Harry{
		MakeArgs: makeArgs,
	}
}

func (harry *Harry) MakeMyDay() {
	for {
		shouldMake, err := harry.watchForRemake()
		if err != nil {
			watchFailedColor.Printf("failed to watch (%s); waiting 1s...\n", err)
			time.Sleep(time.Second)
			continue
		}

		if shouldMake {
			harry.remake()
		}

		harry.wait()
	}
}

func (harry *Harry) watchForRemake() (bool, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return false, fmt.Errorf("failed to construct watcher: %s", err)
	}

	harry.watcher = watcher

	err = watcher.Add("Makefile")
	if err != nil {
		return false, fmt.Errorf("failed to watch Makefile for changes: %s", err)
	}

	detectArgs := []string{"--dry-run", "--debug=v", "--no-builtin-rules"}
	detect := exec.Command("make", append(detectArgs, harry.MakeArgs...)...)

	detectOut, err := detect.StdoutPipe()
	if err != nil {
		return false, fmt.Errorf("failed to construct stdout pipe: %s", err)
	}

	outBuf := bufio.NewReader(detectOut)

	err = detect.Start()
	if err != nil {
		return false, fmt.Errorf("failed spawning 'make' to determine prerequisites: %s", err)
	}

	shouldMake, err := harry.watchPrereqs(outBuf)
	if err != nil {
		return false, fmt.Errorf("failed to watch prerequisites: %s", err)
	}

	err = detect.Wait()
	if err != nil {
		return false, fmt.Errorf("failed to determine prerequisites: %s", err)
	}

	return shouldMake, nil
}

func (harry *Harry) remake() {
	makingColor.Printf("running make %s\n", strings.Join(harry.MakeArgs, " "))

	remake := exec.Command("make", harry.MakeArgs...)
	remake.Stdout = os.Stdout
	remake.Stderr = os.Stderr

	err := remake.Run()
	if err != nil {
		makeFailedColor.Println("make failed")
	} else {
		makeSucceededColor.Println("make succeeded")
	}
}

func (harry *Harry) wait() {
	defer harry.watcher.Close()

	for {
		select {
		case <-harry.watcher.Events:
			return
		case err := <-harry.watcher.Errors:
			watchFailedColor.Printf("encountered error while watching: %s\n", err)
		}
	}
}

func (harry *Harry) watchPrereqs(out *bufio.Reader) (bool, error) {
	var shouldMake bool

	watched := map[string]bool{}

	for {
		line, err := out.ReadString('\n')
		if err == io.EOF {
			return shouldMake, nil
		}

		if err != nil {
			return false, err
		}

		matches := prereqDebugLine.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		prereq := matches[1]

		if matches[2] == "newer" {
			changeColor.Printf("detected change in '%s'\n", prereq)
			shouldMake = true
		}

		parentDir := filepath.Dir(prereq)

		err = harry.watch(watched, parentDir)
		if err != nil {
			return false, err
		}
	}
}

func (harry *Harry) watch(watched map[string]bool, path string) error {
	if watched[path] {
		return nil
	}

	err := harry.watcher.Add(path)
	if err != nil {
		return err
	}

	watched[path] = true

	return nil
}
