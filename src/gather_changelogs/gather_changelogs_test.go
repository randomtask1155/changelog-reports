package main 

import (
	"fmt"
	"testing"
	"os"
)

func TestCheckEnv(t *testing.T) {
	k := "user1"
	os.Setenv("PGHOST", "")
	
	checkEnv(&k, "PGHOST")
	if k != "user1" {
		t.Fatalf("checkEnv Changed user1 to %s when PGHOST is not set\n", k)
	}
	
	os.Setenv("PGHOST", "user2")
	checkEnv(&k, "PGHOST")
	if k != "user2" {
		t.Fatalf("checkEnv did not change user1 to user2 and user currently is %s\n", k)
	}
}

func TestMarshalStruct(t *testing.T) {
	type TestJSON struct {
		ABC string 
		DEF int
	}
	j := &TestJSON{"123", 456}
	b := marshalStruct(j)
	matchString := "{\"ABC\":\"123\",\"DEF\":456}"
	if fmt.Sprintf("%s",b) != matchString {
		t.Fatalf("json marshal failure: %s != %s", b, matchString)
	}
}

func TestCleanCustomerName(t *testing.T) {
	bs := "customer123 % name with!!* Spaces "
	gs := "customer123namewithspaces"
	
	cleanCustomerName(&bs)
	if bs != gs {
		t.Fatalf("strings do not match (bad != good) %s != %s", bs, gs)
	}
}