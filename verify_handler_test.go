package address_verifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyHandler(t *testing.T) {
	t.Run("Something goes in and something goes out", func(t *testing.T) {
		in := make(chan string, 10)
		out := make(chan string, 10)
		handler := NewVerifyHandler(in, out)

		in <- "My String"
		close(in)

		handler.Handle()

		close(out)

		assert.Equal(t, <-out, "My String")
	})
}