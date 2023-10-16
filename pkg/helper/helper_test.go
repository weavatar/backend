package helper

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HelperTestSuite struct {
	suite.Suite
}

func TestHelperTestSuite(t *testing.T) {
	suite.Run(t, &HelperTestSuite{})
}

func (s *HelperTestSuite) TestEmpty() {
	s.True(Empty(""))
	s.True(Empty(nil))
	s.True(Empty([]string{}))
	s.True(Empty(map[string]string{}))
	s.True(Empty(0))
	s.True(Empty(0.0))
	s.True(Empty(false))

	s.False(Empty(" "))
	s.False(Empty([]string{"Goravel"}))
	s.False(Empty(map[string]string{"Goravel": "Goravel"}))
	s.False(Empty(1))
	s.False(Empty(1.0))
	s.False(Empty(true))
}

func (s *HelperTestSuite) TestFirstElement() {
	s.Equal("HaoZi", FirstElement([]string{"HaoZi"}))
}

func (s *HelperTestSuite) TestRandomNumber() {
	s.Len(RandomNumber(10), 10)
}

func (s *HelperTestSuite) TestRandomString() {
	s.Len(RandomString(10), 10)
}

func (s *HelperTestSuite) TestMD5() {
	s.Equal("e10adc3949ba59abbe56e057f20f883e", MD5("123456"))
}

func (s *HelperTestSuite) TestIsURL() {
	s.True(IsURL("https://weavatar.com"))
	s.True(IsURL("https://weavatar.com?name=HaoZi"))
	s.True(IsURL("http://weavatar.com"))
	s.True(IsURL("http://www.weavatar.com"))
	s.False(IsURL("www.weavatar.com"))
	s.False(IsURL("weavatar.com"))
	s.False(IsURL("weavatar.com?name=HaoZi"))

	s.False(IsURL("weavatar"))
}
