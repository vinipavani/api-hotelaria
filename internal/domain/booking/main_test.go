package booking

import (
	"api-hotelaria/internal/database/testdb"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(testdb.SetupIntegrationTests(m))
}
