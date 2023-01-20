package measurementapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	measurementdb "github.com/hieutrtr/influxdb-poc/services/measurements/db"
	"github.com/hieutrtr/influxdb-poc/services/measurements/mocks"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandleListMeasurements(t *testing.T) {
	t.Logf("TestHandleListMeasurements")
	testCases := []struct {
		name          string
		buildStub     func(store *mocks.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().ListMeasurements(gomock.Any(), int64(10), int64(0)).Times(1).Return([]measurementdb.Measurement{
					{
						ID:   primitive.NewObjectID(),
						Name: "test",
					},
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
			},
		},
		{
			name: "internal error",
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().ListMeasurements(gomock.Any(), int64(10), int64(0)).Times(1).Return(nil, measurementdb.NewDBError("internal error", 500))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 500, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocks.NewMockStore(ctrl)
			tc.buildStub(store)
			s := NewServer(store)
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/measurements", nil)
			s.engine.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleCreateMeasurement(t *testing.T) {
	t.Logf("TestHandleCreateMeasurement")
	testCases := []struct {
		name          string
		meaName       string
		buildStub     func(store *mocks.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "success",
			meaName: "test",
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().CreateMeasurement(gomock.Any(), measurementdb.Measurement{
					Name:        "test",
					Description: "test description",
				}).Times(1).Return(&measurementdb.MeasurementID{ID: primitive.NewObjectID()}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
			},
		},
		{
			name:    "internal error",
			meaName: "test",
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().CreateMeasurement(gomock.Any(), measurementdb.Measurement{
					Name:        "test",
					Description: "test description",
				}).Times(1).Return(nil, measurementdb.NewDBError("internal error", 500))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 500, recorder.Code)
			},
		},
		{
			name: "bad request",
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().CreateMeasurement(gomock.Any(), measurementdb.Measurement{
					Description: "test description",
				}).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 400, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocks.NewMockStore(ctrl)
			tc.buildStub(store)
			s := NewServer(store)
			recorder := httptest.NewRecorder()
			data := createMeasurementRequest{
				Name:        tc.meaName,
				Description: "test description",
			}
			body, _ := json.Marshal(data)
			req := httptest.NewRequest("POST", "/measurements", bytes.NewBuffer(body))
			s.engine.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleGetMeasurement(t *testing.T) {
	t.Logf("TestHandleGetMeasurement")
	testCases := []struct {
		name          string
		meaID         string
		buildStub     func(store *mocks.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "success",
			meaID: primitive.NewObjectID().Hex(),
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().GetMeasurement(gomock.Any(), gomock.Any()).Times(1).Return(&measurementdb.Measurement{
					ID:          primitive.NewObjectID(),
					Name:        "test",
					Description: "test description",
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 200, recorder.Code)
			},
		},
		{
			name:  "internal error",
			meaID: primitive.NewObjectID().Hex(),
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().GetMeasurement(gomock.Any(), gomock.Any()).Times(1).Return(nil, measurementdb.NewDBError("internal error", 500))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 500, recorder.Code)
			},
		},
		{
			name:  "bad request",
			meaID: "invalid_id",
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().GetMeasurement(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 400, recorder.Code)
			},
		},
		{
			name:  "not found",
			meaID: primitive.NewObjectID().Hex(),
			buildStub: func(store *mocks.MockStore) {
				store.EXPECT().GetMeasurement(gomock.Any(), gomock.Any()).Times(1).Return(nil, measurementdb.NewErrMeasurementNotFound("not found", 404))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 404, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocks.NewMockStore(ctrl)
			tc.buildStub(store)
			s := NewServer(store)
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/measurements/%s", tc.meaID), nil)
			s.engine.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
