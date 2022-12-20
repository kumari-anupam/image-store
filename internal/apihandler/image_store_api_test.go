package apihandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"githum.com/anupam111/image-store/internal/controller"
	"githum.com/anupam111/image-store/internal/db/dbmodels"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	errFake = errors.New("error")
)

func setupTestEnv(t *testing.T) (
	*gin.Engine,
	*controller.MockImageStore,
	*APIHandler,
) {
	t.Helper()
	log := log.New()
	mockCtrl := gomock.NewController(t)
	mockController := controller.NewMockImageStore(mockCtrl)
	router := gin.New()

	return router, mockController, NewAPIHandler(log, mockController)
}

func Test_CreateImageAlbum(t *testing.T) {
	t.Parallel()

	inputPayload := dbmodels.Album{
		AlbumName: "test-album",
	}

	tests := []struct {
		name    string
		url     string
		payload *dbmodels.Album
		prepare func(
			subs *controller.MockImageStore,
		)
		statusCode int
	}{
		{
			name:    "success",
			url:     "/album",
			payload: &inputPayload,
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().CreateImageAlbum(inputPayload).Return(nil)
			},
			statusCode: 201,
		},
		{
			name:    "internal_server_error",
			url:     "/album",
			payload: &inputPayload,
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().CreateImageAlbum(inputPayload).Return(errFake)
			},
			statusCode: 500,
		},
		{
			name:       "bad_request",
			url:        "/album",
			statusCode: 400,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router, controller, apiHandler := setupTestEnv(t)

			if tt.prepare != nil {
				tt.prepare(controller)
			}
			var payload io.Reader
			if tt.payload != nil {
				temp, _ := json.Marshal(tt.payload)
				payload = bytes.NewReader(temp)
			}
			router.POST("/album", apiHandler.CreateImageAlbum)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, tt.url, payload)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func Test_CreateImage(t *testing.T) {
	t.Parallel()

	inputPayload := dbmodels.Image{
		AlbumName: "test-album",
		ImageName: "test-image",
	}

	tests := []struct {
		name    string
		url     string
		payload *dbmodels.Image
		prepare func(
			subs *controller.MockImageStore,
		)
		statusCode int
	}{
		{
			name:    "success",
			url:     "/image",
			payload: &inputPayload,
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().CreateImage(inputPayload).Return(nil)
			},
			statusCode: 201,
		},
		{
			name:    "internal_server_error",
			url:     "/image",
			payload: &inputPayload,
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().CreateImage(inputPayload).Return(errFake)
			},
			statusCode: 500,
		},
		{
			name:       "bad_request",
			url:        "/image",
			statusCode: 400,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router, controller, apiHandler := setupTestEnv(t)

			if tt.prepare != nil {
				tt.prepare(controller)
			}
			var payload io.Reader
			if tt.payload != nil {
				temp, _ := json.Marshal(tt.payload)
				payload = bytes.NewReader(temp)
			}
			router.POST("/image", apiHandler.CreateImage)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, tt.url, payload)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func Test_DeleteImagesAlbum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		prepare func(
			subs *controller.MockImageStore,
		)
		statusCode int
	}{
		{
			name: "success",
			url:  "/album/test-album",
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().DeleteImageAlbum("test-album").Return(nil)
			},
			statusCode: 204,
		},
		{
			name: "internal_server_error",
			url:  "/album/test-album",
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().DeleteImageAlbum("test-album").Return(errFake)
			},
			statusCode: 500,
		},
		{
			name:       "bad_request",
			url:        "/album",
			statusCode: 404,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router, controller, apiHandler := setupTestEnv(t)

			if tt.prepare != nil {
				tt.prepare(controller)
			}

			var payload io.Reader

			router.DELETE("/album/:albumName", apiHandler.DeleteImageAlbum)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, tt.url, payload)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func Test_DeleteImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		prepare func(
			subs *controller.MockImageStore,
		)
		statusCode int
	}{
		{
			name: "success",
			url:  "/album/images/test-image?albumName=test-album",
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().DeleteImage("test-image", "test-album").Return(nil)
			},
			statusCode: 204,
		},
		{
			name: "internal_server_error",
			url:  "/album/images/test-image?albumName=test-album",
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().DeleteImage("test-image", "test-album").Return(errFake)
			},
			statusCode: 500,
		},
		{
			name:       "bad_request",
			url:        "album/images/",
			statusCode: 404,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router, controller, apiHandler := setupTestEnv(t)

			if tt.prepare != nil {
				tt.prepare(controller)
			}

			var payload io.Reader

			router.DELETE("/album/images/:imageName", apiHandler.DeleteImage)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, tt.url, payload)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func Test_GetImageByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		prepare func(
			subs *controller.MockImageStore,
		)
		statusCode int
	}{
		{
			name: "success",
			url:  "/image/test-image",
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().GetImage("test-image").Return(dbmodels.Image{}, nil)
			},
			statusCode: 200,
		},
		{
			name: "internal_server_error",
			url:  "/image/test-image",
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().GetImage("test-image").Return(dbmodels.Image{}, errFake)
			},
			statusCode: 500,
		},
		{
			name:       "bad_request",
			url:        "/image/",
			statusCode: 404,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router, controller, apiHandler := setupTestEnv(t)

			if tt.prepare != nil {
				tt.prepare(controller)
			}

			var payload io.Reader

			router.GET("/image/:imageName", apiHandler.GetImageByID)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.url, payload)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func Test_GetImages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		prepare func(
			subs *controller.MockImageStore,
		)
		statusCode int
	}{
		{
			name: "success",
			url:  "/album/images?albumName=test-album",
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().GetAllImages("test-album").Return([]dbmodels.Image{}, nil)
			},
			statusCode: 200,
		},
		{
			name: "internal_server_error",
			url:  "/album/images?albumName=test-album",
			prepare: func(subs *controller.MockImageStore) {
				subs.EXPECT().GetAllImages("test-album").Return(nil, errFake)
			},
			statusCode: 500,
		},
		{
			name:       "bad_request",
			url:        "/album/images",
			statusCode: 400,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router, controller, apiHandler := setupTestEnv(t)

			if tt.prepare != nil {
				tt.prepare(controller)
			}

			var payload io.Reader

			router.GET("/album/images", apiHandler.GetAlbumImages)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.url, payload)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
