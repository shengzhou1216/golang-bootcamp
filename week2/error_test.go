package week2 

import (
	"errors"
	"testing"
)

type errorString string

func (e errorString) Error() string {
	return string(e)
}

func New(text string) error {
	return errorString(text)	
}

var ErrNamedType = New("EOF")
var ErrStructType = errors.New("EOF")

func TestError(t *testing.T)   {
	t.Run("测试error相等",func(t *testing.T) {
		if ErrNamedType == New("EOF") {
			t.Log("Named Type Error")
		}
		if ErrStructType == errors.New("EOF") {
			t.Log("Struct Type Error")
		}
	})
}
