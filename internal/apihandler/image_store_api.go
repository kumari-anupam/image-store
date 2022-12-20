package apihandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"githum.com/anupam111/image-store/internal/controller"
	"githum.com/anupam111/image-store/internal/db/dbhandler"
	"githum.com/anupam111/image-store/internal/db/dbmodels"
	"githum.com/anupam111/image-store/internal/models"
	"net/http"
)

// APIHandler handles api.
type APIHandler struct {
	log        *log.Logger
	imageStore controller.ImageStore
}

// NewAPIHandler implements APIHandler.
func NewAPIHandler(logger *log.Logger, imageStore controller.ImageStore) *APIHandler {
	return &APIHandler{
		log:        logger,
		imageStore: imageStore,
	}
}

func (a *APIHandler) CreateImageAlbum(ginCtx *gin.Context) {
	var album dbmodels.Album
	if err := ginCtx.BindJSON(&album); err != nil {
		ginCtx.JSON(http.StatusBadRequest, models.ResponseError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      "BAD-REQUEST",
			MessageDetails: err.Error(),
		})

		return
	}

	a.log.Debugf("album post request payload got: %+v", album)

	err := a.imageStore.CreateImageAlbum(album)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, models.ResponseError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      "INTERNAL-SERVER-ERROR",
			MessageDetails: err.Error(),
		})

		return
	}

	ginCtx.JSON(http.StatusCreated, gin.H{})
}

func (a *APIHandler) CreateImage(ginCtx *gin.Context) {
	var imageModel dbmodels.Image
	if err := ginCtx.BindJSON(&imageModel); err != nil {
		ginCtx.JSON(http.StatusBadRequest, models.ResponseError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      "BAD-REQUEST",
			MessageDetails: err.Error(),
		})

		return
	}

	a.log.Debugf("image post request payload got: %+v", imageModel)

	err := a.imageStore.CreateImage(imageModel)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, models.ResponseError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      "INTERNAL-SERVER-ERROR",
			MessageDetails: err.Error(),
		})

		return
	}

	ginCtx.JSON(http.StatusCreated, gin.H{})
}

func (a *APIHandler) DeleteImageAlbum(ginCtx *gin.Context) {
	albumName := ginCtx.Param("albumName")
	if albumName == "" {
		ginCtx.JSON(http.StatusBadRequest, models.ResponseError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      "BAD-REQUEST",
			MessageDetails: "albumName is empty in query parameter",
		})

		return
	}

	err := a.imageStore.DeleteImageAlbum(albumName)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, models.ResponseError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      "INTERNAL-SERVER-ERROR",
			MessageDetails: err.Error(),
		})

		return
	}

	ginCtx.JSON(http.StatusNoContent, gin.H{})
}

func (a *APIHandler) DeleteImage(ginCtx *gin.Context) {
	imageName := ginCtx.Param("imageName")
	if imageName == "" {
		ginCtx.JSON(http.StatusBadRequest, models.ResponseError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      "BAD-REQUEST",
			MessageDetails: "imageName is empty in request url",
		})

		return
	}

	albumName := ginCtx.Query("albumName")
	if albumName == "" {
		ginCtx.JSON(http.StatusBadRequest, models.ResponseError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      "BAD-REQUEST",
			MessageDetails: "albumName is empty in query parameter",
		})

		return
	}

	err := a.imageStore.DeleteImage(imageName, albumName)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, models.ResponseError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      "INTERNAL-SERVER-ERROR",
			MessageDetails: err.Error(),
		})

		return
	}

	ginCtx.JSON(http.StatusNoContent, gin.H{})
}

func (a *APIHandler) GetImageByID(ginCtx *gin.Context) {
	imageName := ginCtx.Param("imageName")
	if imageName == "" {
		ginCtx.JSON(http.StatusBadRequest, models.ResponseError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      "BAD-REQUEST",
			MessageDetails: "imageName is empty in request url",
		})

		return
	}

	image, err := a.imageStore.GetImage(imageName)
	if err != nil {
		if errors.Is(err, dbhandler.ErrNoDataFound) {
			ginCtx.JSON(http.StatusOK, image)

			return
		}

		ginCtx.JSON(http.StatusInternalServerError, models.ResponseError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      "INTERNAL-SERVER-ERROR",
			MessageDetails: err.Error(),
		})

		return
	}

	ginCtx.JSON(http.StatusOK, image)
}

func (a *APIHandler) GetAlbumImages(ginCtx *gin.Context) {
	albumName := ginCtx.Query("albumName")
	if albumName == "" {
		ginCtx.JSON(http.StatusBadRequest, models.ResponseError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      "BAD-REQUEST",
			MessageDetails: "albumName is empty in request url",
		})

		return
	}

	images, err := a.imageStore.GetAllImages(albumName)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, models.ResponseError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      "INTERNAL-SERVER-ERROR",
			MessageDetails: err.Error(),
		})

		return
	}

	ginCtx.JSON(http.StatusOK, images)
}
