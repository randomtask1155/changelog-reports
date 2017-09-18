package main 

import (
  "testing"
)

func TestChangelogRegexp(t *testing.T) {
	// test exitStatusRegexp
	
}

func TestChangelogNameRegexp(t *testing.T) {
  log := "79_79_1504722951_1504725580_changelog.txt"
  if !changelogNameRegexp.MatchString(log) {
    t.Fatalf("%s does not match changelogNameRegexp", log)
  }
}

func TestExitStatusRegexp(t *testing.T) {
  fail := "Exited with 1.\n"
  pass := "Exited with 0.\n"
  
  if !exitStatusRegexp.MatchString(fail) {
    t.Fatalf("\"%s\" did not match exitStatusRegexp", fail)
  }
  if !exitStatusRegexp.MatchString(fail) {
    t.Fatalf("\"%s\" did not match exitStatusRegexp", pass)
  }
  
  m := exitStatusRegexp.FindStringSubmatch(fail)
  if len(m) != 2 || m[1] != "1" {
    t.Fatalf("exitStatusRegexp did not match status code:%d: matches=%v", len(m),m)
  }
}

func TestStepRegexp(t *testing.T) {
  s := "{\"type\": \"step_started\", \"id\": \"errands.running.p-push-notifications-f04b26227325db529071.push-push-notifications\"}\n"
  
  m := stepRegexp.FindStringSubmatch(s)
  if len(m) != 3 || m[1] != "step_started" || m[2] != "errands.running.p-push-notifications-f04b26227325db529071.push-push-notifications" {
    t.Fatalf("stepRegexp did not match fields:%d: %v", len(m),m)
  }
}

/*
Release info
------------
Name:    push-apps-manager-release
Version: 660.7.9
*/
func TestReleaseInfoRegexp(t *testing.T) {
  s := "Release info\n"
  if !releaseInfoRegexp.MatchString(s) {
    t.Fatalf("%s: did not match releaseInfoRegexp", s)
  }
}

func TestReleaseNameRegexp(t *testing.T) {
  s := "Name:    push-apps-manager-release\n"
  m := releaseNameRegexp.FindStringSubmatch(s)
  if len(m) != 2 || m[1] != "push-apps-manager-release" {
    t.Fatalf("releaseNameRegexp did not match:%d: %v", len(m), m)
  }
}

func TestReleaseVersionRegexp(t *testing.T) {
  s := "Version: 660.7.9\n"
  m := releaseVersionRegexp.FindStringSubmatch(s)
  if len(m) != 2 || m[1] != "660.7.9" {
    t.Fatalf("releaseVersionRegexp did not match:%d: %v", len(m), m)
  }
}



