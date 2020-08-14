package input_test

import (
	"bytes"
	"testing"

	"github.com/jcfox412/logarhythms/internal/input"
	// "github.com/stretchr/testify/assert"
)

func TestRunReturnsPasswordInput(t *testing.T) {
	var stdin bytes.Buffer

	// stdin.Write([]byte("5\n"))

	_ = input.UserInput{
		Reader: &stdin,
	}

	// userInput.PrintMainMenu()
	// assert.NoError(t, err)
	// assert.Equal(t, "hunter2", result)
}
