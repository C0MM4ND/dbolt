package dbolt

import (
	"testing"
)

func TestTx_array_allocatePageStats(t *testing.T) {
	f := newTestArrayFreelist()
	ids := []pgid{2, 3}
	f.readIDs(ids)

	tx := &Tx{
		db: &DB{
			freelist: f,
			pageSize: defaultPageSize,
		},
		meta:  &meta{},
		pages: make(map[pgid]*page),
	}

	prePageCnt := tx.Stats().PageCount
	allocateCnt := f.free_count()

	if _, err := tx.allocate(allocateCnt); err != nil {
		t.Fatal(err)
	}

	if tx.Stats().PageCount != prePageCnt+allocateCnt {
		t.Errorf("Allocated %d but got %d page in stats", allocateCnt, tx.Stats().PageCount)
	}
}

func TestTx_map_allocatePageStats(t *testing.T) {
	f := newTestMapFreelist()
	ids := []pgid{2, 3}
	f.readIDs(ids)

	tx := &Tx{
		db: &DB{
			freelist: f,
			pageSize: defaultPageSize,
		},
		meta:  &meta{},
		pages: make(map[pgid]*page),
	}

	prePageCnt := tx.Stats().PageCount
	allocateCnt := f.free_count()

	if _, err := tx.allocate(allocateCnt); err != nil {
		t.Fatal(err)
	}

	if tx.Stats().PageCount != prePageCnt+allocateCnt {
		t.Errorf("Allocated %d but got %d page in stats", allocateCnt, tx.Stats().PageCount)
	}
}
