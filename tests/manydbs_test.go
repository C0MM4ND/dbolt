package dbolt_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	dbolt "github.com/c0mm4nd/dbolt"
)

func createDb(t *testing.T) (*dbolt.DB, func()) {
	// First, create a temporary directory to be used for the duration of
	// this test.
	tempDirName, err := ioutil.TempDir("", "dboltmemtest")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	path := filepath.Join(tempDirName, "testdb.db")

	bdb, err := dbolt.Open(path, 0600, nil)
	if err != nil {
		t.Fatalf("error creating dbolt db: %v", err)
	}

	cleanup := func() {
		bdb.Close()
		os.RemoveAll(tempDirName)
	}

	return bdb, cleanup
}

func createAndPutKeys(t *testing.T) {
	t.Parallel()

	db, cleanup := createDb(t)
	defer cleanup()

	bucketName := []byte("bucket")

	for i := 0; i < 100; i++ {
		err := db.Update(func(tx *dbolt.Tx) error {
			nodes, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				return err
			}

			var key [16]byte
			rand.Read(key[:])
			if err := nodes.Put(key[:], nil); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestManyDBs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("%d", i), createAndPutKeys)
	}
}
