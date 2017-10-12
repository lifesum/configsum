package client

import "testing"

func TestInmemRepoLookup(t *testing.T) {
	testRepoLookup(t, prepareInmemRepo)
}

func TestInmemRepoLookupNotFound(t *testing.T) {
	testRepoLookupNotFound(t, prepareInmemRepo)
}

func TestInmemTokenRepoLookup(t *testing.T) {
	testTokenRepoLookup(t, prepareInmemTokenRepo)
}

func TestInmemTokenRepoLookupNotFound(t *testing.T) {
	testTokenRepoLookupNotFound(t, prepareInmemTokenRepo)
}

func prepareInmemRepo(t *testing.T) Repo {
	return NewInmemRepo()
}

func prepareInmemTokenRepo(t *testing.T) TokenRepo {
	return NewInmemTokenRepo()
}
