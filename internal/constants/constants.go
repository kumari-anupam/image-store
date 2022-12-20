package constants

const (
	GetImagesQuery                        = `SELECT * FROM Image WHERE "albumName"=$1`
	GetImageByIDQuery                     = `SELECT * FROM Image WHERE "imageName"=$1`
	DeleteImagesOfAlbumQuery              = `DELETE FROM Image WHERE "albumName"=$1`
	DeleteImageWithImageNameAndAlbumQuery = `DELETE FROM Image WHERE "imageName"=$1 AND "albumName"=$2`
	DeleteAlbum                           = `DELETE FROM Album WHERE "albumName"=$1`
)
