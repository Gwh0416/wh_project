package model

const (
	Normal         = 1
	Personal int32 = 1
)

const AESKey = "qwertyuiopasdfghjklzxcvb"

const (
	NoDeleted = 0
	Deleted   = 1
)

const (
	NoArchive = 0
	Archive   = 1
)

const (
	Open    = 0
	Private = 1
	Custom  = 2
)

const (
	Default = "default"
	Simple  = "simple"
)

const (
	NoCollected = 0
	Collected   = 1
)

const (
	NoExecutor = iota
	Executor
)
const (
	NoOwner = iota
	Owner
)
const (
	NoCanRead = iota
	CanRead
)

const (
	NoComment = iota
	Comment
)
