package dbolt

// The largest step that can be taken when remapping the mmap.
const maxMmapStep = 1 << 30 // 1GB

// The data file format version.
const version = 1

// Represents a marker value to indicate that a file is a DBolt DB.
const magic uint32 = 0xCAFEBABE // changed from reversed deadcode to cafebabe

const pgidNoFreelist pgid = 0xffffffffffffffff
