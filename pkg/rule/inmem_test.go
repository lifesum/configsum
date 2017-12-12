package rule

import "testing"

func TestInmemRepoGetByIDNotFound(t *testing.T) {
	testRepoGetByIDNotFound(t, prepareInmemRepo)
}

func TestInmemRepoCreateDuplicate(t *testing.T) {
	testRepoCreateDuplicate(t, prepareInmemRepo)
}

func TestInmemRepoGet(t *testing.T) {
	testRepoGet(t, prepareInmemRepo)
}

func TestInmemRepoUpdateWith(t *testing.T) {
	testRepoUpdateWith(t, prepareInmemRepo)
}

func TestInmemRepoListAll(t *testing.T) {
	testRepoListAll(t, prepareInmemRepo)
}

func TestInmemRepoListActive(t *testing.T) {
	testRepoListActive(t, prepareInmemRepo)
}

func TestInmemRepoListActiveEmpty(t *testing.T) {
	testRepoListActiveEmpty(t, prepareInmemRepo)
}

func TestInmemRepoListAllEmpty(t *testing.T) {
	testRepoListAllEmpty(t, prepareInmemRepo)
}

func TestInmemRepoListDeleted(t *testing.T) {
	testRepoListDeleted(t, prepareInmemRepo)
}

func TestInmemRepoCreateRollout(t *testing.T) {
	testRepoCreateRollout(t, prepareInmemRepo)
}

func prepareInmemRepo(t *testing.T) Repo {
	return NewInmemRuleRepo()
}
