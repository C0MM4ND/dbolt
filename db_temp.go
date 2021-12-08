package dbolt

import (
	"io/ioutil"
	"os"
	"sync"
)

func openTempFile(pattern string, _ int, _ os.FileMode) (*os.File, error) {
	return ioutil.TempFile("", pattern)
}

// Open creates and opens a database at the temporary folder.
// If the file does not exist then it will be created automatically.
// Passing in nil options will cause Bolt to open the database with the default options.
func OpenTemp(pattern string, options *Options) (*DB, error) {
	db := &DB{
		opened:  true,
		Options: options,
	}

	// Set default options if no options are provided.
	if db.Options == nil {
		db.Options = DefaultOptions
	}

	flag := os.O_RDWR // TODO: support multi-file storage
	if db.Options.ReadOnly {
		flag = os.O_RDONLY
		db.readOnly = true
	}

	db.openFile = db.Options.OpenFile
	if db.openFile == nil {
		db.openFile = openTempFile
	}

	// Open data file and separate sync handler for metadata writes.
	var err error
	if db.file, err = db.openFile(pattern, flag|os.O_CREATE, 0); err != nil {
		_ = db.close()
		return nil, err
	}
	db.path = db.file.Name()

	// Lock file so that other processes using Bolt in read-write mode cannot
	// use the database  at the same time. This would cause corruption since
	// the two processes would write meta pages and free pages separately.
	// The database file is locked exclusively (only one process can grab the lock)
	// if !options.ReadOnly.
	// The database file is locked using the shared lock (more than one process may
	// hold a lock at the same time) otherwise (options.ReadOnly is set).
	if err := flock(db, !db.readOnly, db.Options.Timeout); err != nil {
		_ = db.close()
		return nil, err
	}

	// Default values for test hooks
	db.ops.writeAt = db.file.WriteAt

	if db.pageSize = db.Options.PageSize; db.pageSize == 0 {
		// Set the default page size to the OS page size.
		db.pageSize = defaultPageSize
	}

	// Initialize the database if it doesn't exist.
	if info, err := db.file.Stat(); err != nil {
		_ = db.close()
		return nil, err
	} else if info.Size() == 0 {
		// Initialize new files with meta pages.
		if err := db.init(); err != nil {
			// clean up file descriptor on initialization fail
			_ = db.close()
			return nil, err
		}
	} else {
		// Read the first meta page to determine the page size.
		var buf [0x1000]byte
		// If we can't read the page size, but can read a page, assume
		// it's the same as the OS or one given -- since that's how the
		// page size was chosen in the first place.
		//
		// If the first page is invalid and this OS uses a different
		// page size than what the database was created with then we
		// are out of luck and cannot access the database.
		//
		// TODO: scan for next page
		if bw, err := db.file.ReadAt(buf[:], 0); err == nil && bw == len(buf) {
			if m := db.pageInBuffer(buf[:], 0).meta(); m.validate() == nil {
				db.pageSize = int(m.pageSize)
			}
		} else {
			_ = db.close()
			return nil, ErrInvalid
		}
	}

	// Initialize page pool.
	db.pagePool = sync.Pool{
		New: func() interface{} {
			return make([]byte, db.pageSize)
		},
	}

	// Memory map the data file.
	if err := db.mmap(db.Options.InitialMmapSize); err != nil {
		_ = db.close()
		return nil, err
	}

	if db.readOnly {
		return db, nil
	}

	db.loadFreelist()

	// Flush freelist when transitioning from no sync to sync so
	// NoFreelistSync unaware boltdb can open the db later.
	if !db.NoFreelistSync && !db.hasSyncedFreelist() {
		tx, err := db.Begin(true)
		if tx != nil {
			err = tx.Commit()
		}
		if err != nil {
			_ = db.close()
			return nil, err
		}
	}

	// Mark the database as opened and return.
	return db, nil
}
