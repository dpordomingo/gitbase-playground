package handler_test

import (
	"testing"

	"github.com/dpordomingo/gitbase-playground/server/handler"

	"github.com/stretchr/testify/suite"
)

func TestStaticSuite(t *testing.T) {
	suite.Run(t, new(StaticSuite))
}

type StaticSuite struct {
	suite.Suite
}

func (s *StaticSuite) TestConstructor() {
	staticHandler := handler.NewStatic("dir", "server-url")
	s.IsType(&handler.Static{}, staticHandler)
}

/*
func NewStatic(dir, serverURL string) *Static {
	return &Static{
		dir: dir,
		options: options{
			ServerURL: serverURL,
		},
	}
}
*/
