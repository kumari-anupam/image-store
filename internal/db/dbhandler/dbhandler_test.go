package dbhandler

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"githum.com/anupam111/image-store/internal/db/dbconnection"
	"githum.com/anupam111/image-store/internal/db/dbmodels"
	"reflect"
	"testing"
)

func getMocks(t *testing.T) (sqlxmock.Sqlmock, *DBHandler, func()) {
	sqldb, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	log := logrus.New()
	dbHandler := NewDBHandler(log,
		&dbconnection.Pool{
			DB: sqldb,
		},
	)
	finish := func() {
		sqldb.Close()
	}
	return mock, dbHandler, finish
}

func TestCreateAlbum(t *testing.T) {
	mock, dbHandler, finish := getMocks(t)
	defer finish()
	album := dbmodels.Album{
		AlbumName: "test-album",
	}

	tests := []struct {
		name      string
		mock      func()
		errString string
		wantErr   bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("(INSERT INTO Album).*").WithArgs(
					"test-album",
				).WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error",
			mock: func() {
				mock.ExpectExec("(INSERT INTO Album).*").WithArgs(
					"test-album",
				).WillReturnError(errors.New("SQLError"))
				mock.ExpectRollback()
			},
			errString: "SQLError",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			tt.mock()
			err := dbHandler.CreateAlbum(album)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errString)
			} else {
				assert.Nil(t, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestCreateImage(t *testing.T) {
	mock, dbHandler, finish := getMocks(t)
	defer finish()
	image := dbmodels.Image{
		ImageName: "test-image",
		AlbumName: "test-album",
		Image:     "abc.jpg",
	}

	tests := []struct {
		name      string
		mock      func()
		errString string
		wantErr   bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("(INSERT INTO Image).*").WithArgs(
					"test-image",
					"test-album",
					"abc.jpg",
				).WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error",
			mock: func() {
				mock.ExpectExec("(INSERT INTO Image).*").WithArgs(
					"test-image",
					"test-album",
					"abc.jpg",
				).WillReturnError(errors.New("SQLError"))
				mock.ExpectRollback()
			},
			errString: "SQLError",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			tt.mock()
			err := dbHandler.CreateImage(image)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errString)
			} else {
				assert.Nil(t, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

/*func TestDeleteAlbum(t *testing.T) {
	mock, dbHandler, finish := getMocks(t)
	defer finish()

	tests := []struct {
		name      string
		mock      func()
		errString string
		wantErr   bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("DELETE FROM Image WHERE \"albumName\"=").WithArgs(
					"test-album",
				).WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM Album WHERE \"albumName\"=").WithArgs(
					"test-album",
				).WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			tt.mock()
			err := dbHandler.DeleteAlbum("test-album")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errString)
			} else {
				assert.Nil(t, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}*/

func TestFindSubscription(t *testing.T) {
	mock, dbHandler, finish := getMocks(t)
	defer finish()

	images := []dbmodels.Image{
		{
			ImageName: "test-image",
			AlbumName: "test-album",
			Image:     "abc.jpg",
		},
	}

	tests := []struct {
		name                  string
		mock                  func()
		expectedSubscriptions []dbmodels.Image
		errString             string
		expectedErr           bool
	}{
		{
			name: "OK",
			mock: func() {
				columns := []string{"imageName", "albumName", "image"}
				mock.ExpectQuery("SELECT \\* FROM Image WHERE \"albumName\"=").WithArgs("test-album").
					WillReturnRows(sqlxmock.NewRows(columns).AddRow(
						"test-image",
						"test-album",
						"abc.jpg",
					))
			},
			expectedSubscriptions: images,
			expectedErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			var err error
			res, err := dbHandler.GetAllImages("test-album")
			if tt.expectedErr {
				assert.Equal(t, tt.errString, err.Error())
			} else {
				assert.Nil(t, err)
				reflect.DeepEqual(tt.expectedSubscriptions, res)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
