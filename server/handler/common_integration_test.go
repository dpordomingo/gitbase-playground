package handler_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/src-d/gitbase-web/server/handler"
	"github.com/src-d/gitbase-web/server/serializer"
	"github.com/src-d/gitbase-web/server/service"
	testingTools "github.com/src-d/gitbase-web/server/testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
	db                 service.SQLDB
	handler            http.Handler
	requestProcessFunc func(db service.SQLDB) handler.RequestProcessFunc
	IsIntegration      bool
}

func (suite *HandlerSuite) SetupSuite() {
	require := suite.Require()

	//db
	db, err := getDB()
	require.Nil(err)
	require.Nil(db.Ping())
	suite.db = db

	// logger
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	// handler
	queryHandler := handler.APIHandlerFunc(suite.requestProcessFunc(db))
	suite.handler = lg.RequestLogger(logger)(queryHandler)
}

func (suite *HandlerSuite) TearDownSuite() {
	suite.db.Close()
}

type appConfig struct {
	DBConn          string `envconfig:"DB_CONNECTION" default:"root@tcp(localhost:3306)/none?maxAllowedPacket=4194304"`
	BblfshServerURL string `envconfig:"BBLFSH_SERVER_URL" default:"127.0.0.1:9432"`
	IsIntegration   bool   `envconfig:"INTEGRATION_TESTS" default:"false"`
}

func getDB() (service.SQLDB, error) {
	var conf appConfig
	envconfig.MustProcess("GITBASEPG", &conf)

	if conf.IsIntegration {
		return sql.Open("mysql", conf.DBConn)
	}

	return &testingTools.MockDB{}, nil
}

func errorResponse(res *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)

	return resBody, err
}

func firstErr(require *require.Assertions, resBody map[string]interface{}) map[string]interface{} {
	require.NotEmpty(resBody["errors"].([]interface{}))
	return resBody["errors"].([]interface{})[0].(map[string]interface{})
}

func firstRow(require *require.Assertions, res *httptest.ResponseRecorder) map[string]interface{} {
	var resBody serializer.Response
	json.Unmarshal(res.Body.Bytes(), &resBody)
	require.NotEmpty(resBody.Data.([]interface{}))
	return resBody.Data.([]interface{})[0].(map[string]interface{})
}

func okResponse(require *require.Assertions, res *httptest.ResponseRecorder) {
	var resBody serializer.Response
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	require.Nil(err, res.Body.String())

	require.Equal(http.StatusOK, res.Code, resBody.Errors)

	require.Equal(res.Code, resBody.Status)
	require.NotEmpty(resBody.Data)
	require.NotEmpty(resBody.Meta)
}

func isIntegration() bool {
	var conf appConfig
	envconfig.MustProcess("GITBASEPG", &conf)

	return conf.IsIntegration
}

func bblfshServerURL() string {
	var conf appConfig
	envconfig.MustProcess("GITBASEPG", &conf)

	return conf.BblfshServerURL
}
