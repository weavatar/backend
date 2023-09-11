package tests

import (
	"github.com/goravel/framework/testing"

	"weavatar/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
