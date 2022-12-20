package controller

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"githum.com/anupam111/image-store/internal/db/dbhandler"
	"githum.com/anupam111/image-store/internal/db/dbmodels"
	"testing"
)

var (
	errFake        = errors.New("error")
	subscriptionID = 1
)

func testSetUp(t *testing.T) (
	*gomock.Controller,
	*dbhandler.MockImageStore,
	*ImageController,
) {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	mockHandler := dbhandler.NewMockImageStore(mockCtrl)
	log := logrus.New()
	return mockCtrl, mockHandler, NewImageController(
		log,
		mockHandler,
	)
}

func TestCreateAlbum(t *testing.T) {
	t.Parallel()

	album := dbmodels.Album{
		AlbumName: "test-album",
	}

	tests := []struct {
		name    string
		input   dbmodels.Album
		prepare func(
			subs *dbhandler.MockImageStore,
		)
		expectedError error
	}{
		{
			name:  "success",
			input: album,
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().CreateAlbum(album).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "internal_server",
			input: album,
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().CreateAlbum(album).Return(errFake)
			},
			expectedError: fmt.Errorf("error while creating image album, %w", errFake),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, mockDbHandler, controller := testSetUp(t)
			if tt.prepare != nil {
				tt.prepare(mockDbHandler)
			}

			err := controller.CreateImageAlbum(tt.input)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestCreateImage(t *testing.T) {
	t.Parallel()

	image := dbmodels.Image{
		AlbumName: "test-album",
		ImageName: "test-image",
	}

	tests := []struct {
		name    string
		input   dbmodels.Image
		prepare func(
			subs *dbhandler.MockImageStore,
		)
		expectedError error
	}{
		{
			name:  "success",
			input: image,
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().CreateImage(image).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "internal_server",
			input: image,
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().CreateImage(image).Return(errFake)
			},
			expectedError: fmt.Errorf("error while creating image, %w", errFake),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, mockDbHandler, controller := testSetUp(t)
			if tt.prepare != nil {
				tt.prepare(mockDbHandler)
			}

			err := controller.CreateImage(tt.input)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestDeleteAlbum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(
			subs *dbhandler.MockImageStore,
		)
		expectedError error
	}{
		{
			name: "success",
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().DeleteAlbum("test-album").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "internal_server",
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().DeleteAlbum("test-album").Return(errFake)
			},
			expectedError: fmt.Errorf("error while deleting image album, %w", errFake),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, mockDbHandler, controller := testSetUp(t)
			if tt.prepare != nil {
				tt.prepare(mockDbHandler)
			}

			err := controller.DeleteImageAlbum("test-album")
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestDeleteImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(
			subs *dbhandler.MockImageStore,
		)
		expectedError error
	}{
		{
			name: "success",
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().DeleteImageWithImageName("test-image", "test-album").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "error",
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().DeleteImageWithImageName("test-image", "test-album").Return(errFake)
			},
			expectedError: fmt.Errorf("error while deleting image, %w", errFake),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, mockDbHandler, controller := testSetUp(t)
			if tt.prepare != nil {
				tt.prepare(mockDbHandler)
			}

			err := controller.DeleteImage("test-image", "test-album")
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestGetImageByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(
			subs *dbhandler.MockImageStore,
		)
		expectedError error
		expectedImage dbmodels.Image
	}{
		{
			name: "success",
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().GetImageByID("test-image").Return(dbmodels.Image{}, nil)
			},
			expectedError: nil,
			expectedImage: dbmodels.Image{},
		},
		{
			name: "error",
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().GetImageByID("test-image").Return(dbmodels.Image{}, errFake)
			},
			expectedError: fmt.Errorf("error while getting image, %w", errFake),
			expectedImage: dbmodels.Image{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, mockDbHandler, controller := testSetUp(t)
			if tt.prepare != nil {
				tt.prepare(mockDbHandler)
			}

			image, err := controller.GetImage("test-image")
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedImage, image)
		})
	}
}

func TestGetAllImages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(
			subs *dbhandler.MockImageStore,
		)
		expectedError error
		expectedImage []dbmodels.Image
	}{
		{
			name: "success",
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().GetAllImages("test-album").Return(nil, nil)
			},
			expectedError: nil,
			expectedImage: nil,
		},
		{
			name: "error",
			prepare: func(
				subs *dbhandler.MockImageStore,
			) {
				subs.EXPECT().GetAllImages("test-album").Return(nil, errFake)
			},
			expectedError: fmt.Errorf("error while getting images with album name = test-album,error :  %w", errFake),
			expectedImage: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, mockDbHandler, controller := testSetUp(t)
			if tt.prepare != nil {
				tt.prepare(mockDbHandler)
			}

			image, err := controller.GetAllImages("test-album")
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedImage, image)
		})
	}
}
