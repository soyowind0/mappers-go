package model

type Mission struct {
	UniqueName       string `gorm:"primaryKey"`
	Command          string
	FileContent      string
	FileName         string
	WorkingDirectory string
	Status           string
	Output           string
}
