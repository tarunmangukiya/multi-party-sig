package xor

import (
	"errors"

	gogo "github.com/gogo/protobuf/types"
	"github.com/taurusgroup/cmp-ecdsa/pkg/message"
	"github.com/taurusgroup/cmp-ecdsa/pkg/party"
	"github.com/taurusgroup/cmp-ecdsa/pkg/round"
	"github.com/taurusgroup/cmp-ecdsa/pkg/types"
)

// Round2 embeds Round1 so that it has access to previous information
type Round2 struct {
	*Round1

	// resultXOR is the XOR of all messages sent by parties
	resultXOR []byte
}

// Result implements round.Final and returns the result
func (r *Round2) Result() interface{} { return Result(r.resultXOR) }

// ProcessMessage casts the content to the appropriate type and stores the content.
func (r *Round2) ProcessMessage(from party.ID, content message.Content) error {
	body := content.(*Round2Message)
	// store the received value
	r.received[from] = body.Value
	return nil
}

// Finalize does not send any messages, but computes the output resulting from the received messages
func (r *Round2) Finalize(chan<- *message.Message) error {
	r.resultXOR = make([]byte, 32)
	for _, received := range r.received {
		for i := range r.resultXOR {
			r.resultXOR[i] ^= received[i]
		}
	}
	return nil
}

// Next returns nil to indicate the protocol is done, since this is the final round
func (r *Round2) Next() round.Round { return nil }

// MessageContent returns an uninitialized Round2Message used to unmarshal contents embedded in message.Message.
func (r *Round2) MessageContent() message.Content { return &Round2Message{} }

// Round2Message is the message sent in Round1 and received in Round2
// It contains the
type Round2Message struct {
	// simple wrapper for bytes
	gogo.BytesValue
}

// Validate returns an error if basic information about the message was incorrect.
// This is a good place to check buffer lengths and non nil values
func (m *Round2Message) Validate() error {
	if len(m.Value) != 32 {
		return errors.New("value should be 32 bytes long")
	}
	return nil
}

// RoundNumber indicates which round this message is supposed to be received in.
func (m *Round2Message) RoundNumber() types.RoundNumber { return 2 }
