package main 

import (
  "regexp"
)

var (
  exitStatusRegexp = regexp.MustCompile(`^Exited\s+with\s+(\d+).$`)
  stepRegexp = regexp.MustCompile(`^{"type":\s+"\S+",\s+"id":\s+"\S+"}$`)
  releaseInfoRegexp = regexp.MustCompile(`^Release\s+info$`)
  releaseNameRegexp = regexp.MustCompile(`^Name:\s+(\S+)$`)
  releaseVersionRegexp = regexp.MustCompile(`^Version:\s+(\S+)$`)
)

// Step is equivalent to opsman json step delimiter {"type": "step_started", "id": "errands.running.p-push-notifications-f04b26227325db529071.push-push-notifications"}
type Step struct {
	Type string `json:"type"`
	ID string `json:"id"`
}

func scanReleaseinfo(){
/*
Release info
------------
Name:    binary-offline-buildpack
Version: 1.0.12
*/
}

func main() {
  
}