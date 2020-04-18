package spool

import (
	"io"

	dgo "github.com/bwmarrin/discordgo"
)

type Head string

const Zero = Head("00000000000000000000")

// spool.T is the structure which encapsulates all state necessary to iterate
// over a paginated channel message history.
type T struct {
	ChID string
	Head
}

// Unroll attempts to iter() on the message history of the given Discord
// channel ID, in descending chronological order of appearance, and blocks
// until some error is encountered, all messages have been read (io.EOF), or
// iter returns false.
func (spl *T) Unroll(s *dgo.Session, iter func(*dgo.Message) bool) (err error) {
	if spl.Head == "" {
		spl.Head = Zero
	}
	msgs := make([]*dgo.Message, 0, 100)
	for {
		msgs, err = s.ChannelMessages(spl.ChID, 100, string(spl.Head), "", "")
		if err != nil {
			return
		}
		if len(msgs) > 0 {
			tail := Head(msgs[len(msgs)-1].ID)
			if head == tail {
				msgs = msgs[:0]
			} else {
				spl.Head = Head(msgs[len(msgs)-1].ID)
			}
		}
		for _, msg := range msgs {
			if !iter(msg) {
				return
			}
		}
		if len(msgs) < 100 {
			err = io.EOF
			return
		}
	}
	return
}
