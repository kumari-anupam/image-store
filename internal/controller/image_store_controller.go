package controller

//go:generate mockgen -source ./image_store_controller.go -package controller -destination image_store_controller_mock.go

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"githum.com/anupam111/image-store/internal/db/dbhandler"
	"githum.com/anupam111/image-store/internal/db/dbmodels"
)

type ImageStore interface {
	CreateImageAlbum(album dbmodels.Album) error
	DeleteImageAlbum(albumName string) error
	CreateImage(image dbmodels.Image) error
	DeleteImage(imageName, albumName string) error
	GetImage(id string) (dbmodels.Image, error)
	GetAllImages(albumName string) ([]dbmodels.Image, error)
}

type ImageController struct {
	log        *log.Logger
	imageStore dbhandler.ImageStore
}

func NewImageController(log *log.Logger, imageStore dbhandler.ImageStore) *ImageController {
	return &ImageController{log: log,
		imageStore: imageStore}
}

func (i *ImageController) CreateImageAlbum(album dbmodels.Album) error {
	err := i.imageStore.CreateAlbum(album)
	if err != nil {
		return fmt.Errorf("error while creating image album, %w", err)
	}

	return nil
}

func (i *ImageController) DeleteImageAlbum(albumName string) error {
	err := i.imageStore.DeleteAlbum(albumName)
	if err != nil {
		return fmt.Errorf("error while deleting image album, %w", err)
	}

	return nil
}

func (i *ImageController) CreateImage(image dbmodels.Image) error {
	err := i.imageStore.CreateImage(image)
	if err != nil {
		return fmt.Errorf("error while creating image, %w", err)
	}

	return nil
}

func (i *ImageController) DeleteImage(imageName, albumName string) error {
	err := i.imageStore.DeleteImageWithImageName(imageName, albumName)
	if err != nil {
		return fmt.Errorf("error while deleting image, %w", err)
	}

	return nil
}

func (i *ImageController) GetImage(id string) (dbmodels.Image, error) {
	image, err := i.imageStore.GetImageByID(id)
	if err != nil {
		return dbmodels.Image{}, fmt.Errorf("error while getting image, %w", err)
	}

	return image, nil
}

func (i *ImageController) GetAllImages(albumName string) ([]dbmodels.Image, error) {
	images, err := i.imageStore.GetAllImages(albumName)
	if err != nil {
		return nil, fmt.Errorf("error while getting images with "+
			"album name = %s,error :  %w", albumName, err)
	}

	return images, nil
}
