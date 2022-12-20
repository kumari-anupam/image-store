package dbmodels

type Album struct {
	AlbumName string `db:"albumName"`
}

type Image struct {
	ImageName string `db:"imageName"`
	AlbumName string `db:"albumName"`
	Image     string `db:"image"`
}
