package storage

type Repository interface {
	Get(string) (string, bool)
	Write(string, string)
}
