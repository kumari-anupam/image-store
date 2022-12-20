package dbconnection

import (
	"githum.com/anupam111/image-store/internal/config"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func getMockController(t *testing.T) (*gomock.Controller, *MockDBConnector) {
	mockCtrl := gomock.NewController(t)
	mockConnector := NewMockDBConnector(mockCtrl)

	return mockCtrl, mockConnector
}

func Test_setConnector(t *testing.T) {
	tests := []struct {
		name     string
		dbConfig *config.DBConfig
		want     *config.DBConfig
	}{
		{
			name:     "Default",
			dbConfig: &config.DBConfig{},
			want: &config.DBConfig{
				Connector: &Connector{},
			},
		},
		{
			name: "Override",
			dbConfig: &config.DBConfig{
				Connector: &MockDBConnector{},
			},
			want: &config.DBConfig{
				Connector: &MockDBConnector{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setConnector(tt.dbConfig)
			assert.Equal(t, tt.dbConfig, tt.want)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		dbConfig *config.DBConfig
		prepare  func(connector *MockDBConnector)
		want     *Pool
	}{
		{
			name: "OK",
			prepare: func(connector *MockDBConnector) {
				connector.EXPECT().Connect(
					"postgres",
					"host=host port=5432 user=dummy password='dummy' dbname=dummy sslmode=disable",
				).Return(&sqlx.DB{})
			},
			dbConfig: &config.DBConfig{
				DriverName: "postgres",
				Host:       "host",
				Port:       5432,
				Name:       "dummy",
				Username:   "dummy",
				Password:   "dummy",
			},
			want: &Pool{
				DB: &sqlx.DB{},
			},
		},
	}

	for _, tt := range tests {
		_, mockConnector := getMockController(t)
		tt.prepare(mockConnector)
		tt.dbConfig.Connector = mockConnector
		got := New(tt.dbConfig)
		assert.Equal(t, got, tt.want)
	}
}
