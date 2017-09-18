package main 

import (
  "regexp"
  "flag"
  "fmt"
  "archive/tar"
  "os"
  "bytes"
  "io"
)

var (
  exitStatusRegexp = regexp.MustCompile(`^Exited\s+with\s+(\d+)\.\n$`)
  stepRegexp = regexp.MustCompile(`^{"type":\s+"(\S+)",\s+"id":\s+"(\S+)"}\n$`)
  releaseInfoRegexp = regexp.MustCompile(`^Release\s+info\n$`)
  releaseNameRegexp = regexp.MustCompile(`^Name:\s+(\S+)\n$`)
  releaseVersionRegexp = regexp.MustCompile(`^Version:\s+(\S+)\n$`)
  changelogNameRegexp = regexp.MustCompile(`\S+_changelog.txt$`)
  
  //Collors
	textColor = "\033[94m"
	errorColor = "\033[91m"
	warnColor = "\033[33m"
	termColor = "\033[0m"
  
  inputlog = flag.String("i", "", "compressed tar file for processing")
  lastStep = "NONE"
)

// Step is equivalent to opsman json step delimiter {"type": "step_started", "id": "errands.running.p-push-notifications-f04b26227325db529071.push-push-notifications"}
type Step struct {
	Type string `json:"type"`
	ID string `json:"id"`
}

/*
Exited with 0.
*/
func scanExitStatus(s string){
  m := exitStatusRegexp.FindStringSubmatch(s)
  if len(m) != 2 {
    logWarn(fmt.Sprintf("could not find return code in exist status string: %s", s))
    return
  }
  if m[1] == "0" {
    logInfo(fmt.Sprintf("Changelog Executed successfully: %s", s))
  } else {
    logInfo(fmt.Sprintf("Changelog Exited on step %s: %s\n", lastStep, s))
  }
}

/*
Release info
------------
Name:    binary-offline-buildpack
Version: 1.0.12
*/
func scanReleaseinfo(buf *bytes.Buffer){
  release := ""
  
  // next line should be "------------" so lets read and dump it
  _, err := buf.ReadString('\n')
  if err != nil {
    logError(fmt.Sprintf("Scanning release line failed: %s", err))
    return
  }
  
  s, err := buf.ReadString('\n')
  if err != nil {
    logError(fmt.Sprintf("Something went wrong while scanning release info name: %s", err))
    return
  }
  
  m := releaseNameRegexp.FindStringSubmatch(s)
  if len(m) != 2 {
    logWarn(fmt.Sprintf("Could not match release name in \"%s\"", s))
    return
  }
  release += fmt.Sprintf("%s", m[1])
  
  s, err = buf.ReadString('\n')
  if err != nil {
    logError(fmt.Sprintf("Something went wrong while scanning release info version: %s", err))
    return
  }
  
  m = releaseVersionRegexp.FindStringSubmatch(s)
  if len(m) != 2 {
    logWarn(fmt.Sprintf("Could not match release version in \"%s\"", s))
    return
  }
  logInfo(fmt.Sprintf("Found Release: %s: %s", release, m[0]))
}

/*
{"type": "step_started", "id": "errands.running.p-push-notifications-f04b26227325db529071.push-push-notifications"}
*/
func scanStep(s string) {
  m := stepRegexp.FindStringSubmatch(s)
  if len(m) != 3 {
    logWarn(fmt.Sprintf("Could not find step properties in string: %s:%d: %v", s, len(m), m))
    return
  }
  lastStep = m[2] // set last step the the ID field
}

/*
TODO
*/
func scanError() {}

func logError(s string){
	fmt.Printf("~# %s%s%s\n", errorColor, s, termColor)
}
func logWarn(s string){
	fmt.Printf("~# %s%s%s\n", warnColor, s, termColor)
}
func logInfo(s string){
	fmt.Printf("~# %s%s%s\n", textColor, s, termColor)
}

func readToBuf(buf *bytes.Buffer, tr *tar.Reader) {
  _, err := io.Copy(buf, tr)
  if err != nil {
    panic(err.Error())
  }
}

func processTar(tf string) error {
  fh, err := os.Open(tf)
  if err != nil {
    return err
  }
  defer fh.Close()
  tr := tar.NewReader(fh)
  
  
  for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
    if err != nil {
			panic(err.Error())
		}
    
    if !changelogNameRegexp.MatchString(hdr.Name) {
      continue
    } 
    logInfo(fmt.Sprintf("Processing %s", hdr.Name))
    
    buf := new(bytes.Buffer)
    readToBuf(buf, tr)
    processLog(buf)
  }
  return nil
}

func processLog(buf *bytes.Buffer) {
  for {
    s, err := buf.ReadString('\n')
    if err == io.EOF {
      break
    }
    if err != nil {
      panic(err.Error())
    }
    
    if exitStatusRegexp.MatchString(s) {
      scanExitStatus(s)
    } else if stepRegexp.MatchString(s) {
      scanStep(s)
    } else if releaseInfoRegexp.MatchString(s) {
      scanReleaseinfo(buf)
    }
  }
}

func main() {
  flag.Parse()
  processTar(*inputlog)  
}