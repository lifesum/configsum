package pg

// Default URIs.
const (
	DefaultDevURI  = "postgres://%s@127.0.0.1:5432/configsum_dev?sslmode=disable&connect_timeout=5"
	DefaultTestURI = "postgres://%s@127.0.0.1:5432/configsum_test?sslmode=disable&connect_timeout=5"
)
