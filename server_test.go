// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"github.com/bmizerany/assert"
	"labix.org/v2/mgo"
	"log"
	"testing"
)

var (
	testServ  Server
	testMongo *mgo.Database
)

func col() *mgo.Collection {
	return testMongo.C("authorizations")
}

func setup(t *testing.T) {
	/*
		if testServ != nil {
			t.Log("Using existing testAuthServer\n")
			return
		}
		t.Log("Initializing testAuthServer\n")
	*/
	log.SetFlags(log.Ltime | log.Lshortfile)
	session, err := mgo.Dial("localhost")
	if err != nil {
		t.Fatal(err)
	}
	testMongo = session.DB("test_btoken")
	err = testMongo.DropDatabase()
	if err != nil {
		t.Fatal(err)
	}
	testServ, err = NewMongoServer(testMongo)
	if err != nil {
		t.Fatal(err)
	}
	return
}

func TestIssueToken(t *testing.T) {
	setup(t)
	user := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	req := AuthRequest{
		User:   user,
		Scopes: scopes,
	}
	token, err := testServ.IssueToken(req)
	if err != nil {
		t.Error(err)
	}
	c := col()
	cnt, err := c.Count()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, cnt)
	query := struct {
		Token string
	}{
		Token: token,
	}
	q := c.Find(query)
	cnt, err = q.Count()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, cnt)
	a := Authorization{}
	err = q.One(&a)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, user, a.User)
	for _, scope := range scopes {
		_, ok := a.Scopes[scope]
		assert.T(t, ok, "Expected scope: ", scope)
	}
}

func TestGetAuthorization(t *testing.T) {
	setup(t)
	user := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	req := AuthRequest{
		User:   user,
		Scopes: scopes,
	}
	token, err := testServ.IssueToken(req)
	if err != nil {
		t.Error(err)
	}
	a, err := testServ.GetAuthorization(token)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, user, a.User)
	for _, scope := range scopes {
		_, ok := a.Scopes[scope]
		assert.T(t, ok, "Expected scope: ", scope)
	}
}
