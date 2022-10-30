package filesys

type FsIface interface {
	CreateDir(name string) error
}
