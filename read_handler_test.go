package addrvrf_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
)

func TestReadHandler(t *testing.T) {
	type fixture struct {
		buffer  *bytes.Buffer
		output  chan *addrvrf.Envelope
		handler *addrvrf.ReadHandler
	}

	setup := func() *fixture {
		buffer := bytes.NewBufferString("")
		output := make(chan *addrvrf.Envelope, 10)
		handler := addrvrf.NewReadHandler(buffer, output)

		writeLine(buffer, "Street,City,State,ZIPCode")

		return &fixture{
			buffer:  buffer,
			output:  output,
			handler: handler,
		}
	}

	t.Run("Read a CSV line", func(t *testing.T) {
		f := setup()

		writeLine(f.buffer, "A1,B1,C1,D1")

		f.handler.Handle()
		close(f.output)

		assertEnvelopeSent(t, addrvrf.InitialSequence, <-f.output)
	})

	t.Run("Read multiple lines and create Envelopes", func(t *testing.T) {
		f := setup()

		writeLine(f.buffer, "A1,B1,C1,D1")
		writeLine(f.buffer, "A2,B2,C2,D2")
		writeLine(f.buffer, "A3,B3,C3,D3")
		writeLine(f.buffer, "A4,B4,C4,D4")
		writeLine(f.buffer, "A5,B5,C5,D5")

		f.handler.Handle()
		close(f.output)

		assertEnvelopeSent(t, addrvrf.InitialSequence, <-f.output)
		assertEnvelopeSent(t, addrvrf.InitialSequence+1, <-f.output)
		assertEnvelopeSent(t, addrvrf.InitialSequence+2, <-f.output)
		assertEnvelopeSent(t, addrvrf.InitialSequence+3, <-f.output)
		assertEnvelopeSent(t, addrvrf.InitialSequence+4, <-f.output)
	})
}

func writeLine(buffer *bytes.Buffer, line string) {
	buffer.WriteString(line + "\n")
}

func assertEnvelopeSent(t *testing.T, seq int, actual *addrvrf.Envelope) {
	num := strconv.Itoa(seq)
	expected := &addrvrf.Envelope{
		Sequence: seq,
		Input: addrvrf.AddressInput{
			Street:  "A" + num,
			City:    "B" + num,
			State:   "C" + num,
			ZIPCode: "D" + num,
		},
	}

	assert.Equal(t, expected, actual)
}
