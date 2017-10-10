package config

import "testing"

func TestInmenUserRepoGetLatest(t *testing.T) {
	testUserRepoGetLatest(t, prepareInmenUserRepo)
}

func TestInmemUserRepoGetLatestNotFound(t *testing.T) {
	testUserRepoGetLatestNotFound(t, prepareInmenUserRepo)
}

func TestInmemUserRepoAppendDuplicate(t *testing.T) {
	testUserRepoAppendDuplicate(t, prepareInmenUserRepo)
}

func prepareInmenUserRepo(t *testing.T) UserRepo {
	r, err := NewInmemUserRepo()
	if err != nil {
		t.Fatal(err)
	}

	return r
}
