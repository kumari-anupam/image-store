package dbhandler

//go:generate mockgen -source ./dbhandler.go -package dbhandler -destination dbhandler_mock.go

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"githum.com/anupam111/image-store/internal/constants"
	"githum.com/anupam111/image-store/internal/db/dbconnection"
	"githum.com/anupam111/image-store/internal/db/dbmodels"
)

var (
	ErrDuplicate   = errors.New("duplicate insertion request")
	ErrNoDataFound = errors.New("no record found")
)

type ImageStore interface {
	CreateAlbum(album dbmodels.Album) error
	CreateImage(image dbmodels.Image) error
	DeleteAlbum(albumName string) error
	DeleteAllImagesOfAlbum(albumName string) error
	DeleteImageWithImageName(imageName, albumName string) error
	GetImageByID(imageID string) (dbmodels.Image, error)
	GetAllImages(albumName string) ([]dbmodels.Image, error)
}

type DBHandler struct {
	log        *log.Logger
	connection *dbconnection.Pool
}

// NewDBHandler implements DBHandler.
func NewDBHandler(log *log.Logger, connection *dbconnection.Pool) *DBHandler {
	return &DBHandler{
		log:        log,
		connection: connection,
	}
}

func (db *DBHandler) CreateAlbum(album dbmodels.Album) error {
	txn := db.connection.DB.MustBegin()
	if _, err := txn.NamedExec(
		`INSERT INTO Album(
			"albumName"
		) VALUES(
			:albumName
		)`,
		album,
	); err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code.Name() == "unique_violation" {
				return ErrDuplicate
			}
		}

		return fmt.Errorf("%w", handlerError(err, txn))
	}

	if err := handlerError(nil, txn); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (db *DBHandler) CreateImage(image dbmodels.Image) error {
	txn := db.connection.DB.MustBegin()
	if _, err := txn.NamedExec(
		`INSERT INTO Image(
			"imageName",
            "albumName",
        	"image"
		) VALUES(
		    :imageName,
			:albumName,
			:image
		)`,
		image,
	); err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code.Name() == "unique_violation" {
				return ErrDuplicate
			}
		}

		return fmt.Errorf("%w", handlerError(err, txn))
	}

	if err := handlerError(nil, txn); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (db *DBHandler) DeleteAlbum(albumName string) error {
	err := db.DeleteAllImagesOfAlbum(albumName)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	tx := db.connection.DB.MustBegin()
	_, err = db.connection.DB.Exec(constants.DeleteAlbum, albumName)

	if err = handlerError(err, tx); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (db *DBHandler) DeleteImageWithImageName(imageName, albumName string) error {
	tx := db.connection.DB.MustBegin()
	_, err := db.connection.DB.Exec(constants.DeleteImageWithImageNameAndAlbumQuery,
		imageName, albumName)

	if err = handlerError(err, tx); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (db *DBHandler) DeleteAllImagesOfAlbum(albumName string) error {
	tx := db.connection.DB.MustBegin()
	_, err := db.connection.DB.Exec(constants.DeleteImagesOfAlbumQuery, albumName)

	if err = handlerError(err, tx); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (db *DBHandler) GetImageByID(imageName string) (dbmodels.Image, error) {
	res := dbmodels.Image{}

	if err := db.connection.DB.Get(&res, constants.GetImageByIDQuery, imageName); err != nil {
		// To handle row not exist
		if errors.Is(err, sql.ErrNoRows) {
			db.log.Errorf("image not found for the request=%s", imageName)

			return dbmodels.Image{}, ErrNoDataFound
		}

		db.log.Errorf("error while getting image from the table for the requeste=%s", imageName)

		return dbmodels.Image{}, fmt.Errorf("%w", err)
	}

	return res, nil
}

func (db *DBHandler) GetAllImages(albumName string) ([]dbmodels.Image, error) {
	images := []dbmodels.Image{}

	rows, err := db.connection.DB.Queryx(constants.GetImagesQuery, albumName)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	for rows.Next() {
		image := dbmodels.Image{}

		err = rows.StructScan(&image)
		if err != nil {
			db.log.Errorf("erro while converting data found from db to image struct: %v", err)

			return nil, fmt.Errorf("error while scanning db data to image: %w", err)
		}

		images = append(images, image)
	}

	defer rows.Close()

	return images, nil
}
