package config

import "testing"

func TestInmemBaseRepoCreateDuplicate(t *testing.T) {
	testBaseRepoCreateDuplicate(t, prepareInmemBaseRepo)
}

func TestInmemBaseRepoGetByID(t *testing.T) {
	testBaseRepoGetByID(t, prepareInmemBaseRepo)
}

func TestInmemBaseRepoGetByIDNotFound(t *testing.T) {
	testBaseRepoGetByIDNotFound(t, prepareInmemBaseRepo)
}

func TestInmemBaseRepoGetByName(t *testing.T) {
	testBaseRepoGetByName(t, prepareInmemBaseRepo)
}

func TestInmemBaseRepoGetByNameNotFound(t *testing.T) {
	testBaseRepoGetByNameNotFound(t, prepareInmemBaseRepo)
}

func TestInmemBaseRepoList(t *testing.T) {
	testBaseRepoList(t, prepareInmemBaseRepo)
}

func TestInmemBaseRepoUpdate(t *testing.T) {
	testBaseRepoUpdate(t, prepareInmemBaseRepo)
}

func TestInmenUserRepoGetLatest(t *testing.T) {
	testUserRepoGetLatest(t, prepareInmenUserRepo)
}

func TestInmemUserRepoGetLatestNotFound(t *testing.T) {
	testUserRepoGetLatestNotFound(t, prepareInmenUserRepo)
}

func TestInmemUserRepoAppendDuplicate(t *testing.T) {
	testUserRepoAppendDuplicate(t, prepareInmenUserRepo)
}

func prepareInmemBaseRepo(t *testing.T) BaseRepo {
	return NewInmemBaseRepo(nil)
}

func prepareInmenUserRepo(t *testing.T) UserRepo {
	return NewInmemUserRepo()
}
