package crawler

import (
	"log"
	"time"

	"github.com/ahmdrz/goinsta"
)

// Message ...
type Message struct {
	Username string
	Time     time.Time
	UserId   int64
	Text     string
}

var _users map[int64]string

// GrabConversations all available
func (s *service) GrabConversations() (map[string][]Message, error) {
	log.SetPrefix("GrabConversations ")
	if err := s.api.Inbox.Sync(); err != nil {
		s.l.Printf("inbox sync error `%v`", err)
		return nil, err
	}

	if err := s.api.Inbox.Sync(); err != nil {
		s.l.Printf("inbox sync error `%v`", err)
		return nil, err
	}

	if _users == nil {
		_users = make(map[int64]string)
	}

	result := make(map[string][]Message)
	for _, c := range s.api.Inbox.Conversations {
		//s.l.Printf("users: %#v", c.Users)
		r := s.recipient(c.Users)
		if r == "" {
			continue
		}
		if len(c.Users) > 2 {
			s.l.Printf("conversation users more then 2")
		}
		for c.Next() {
			for _, i := range c.Items {
				if i.Type == "text" {
					result[r] = append(result[r], Message{
						Username: _users[i.UserID],
						Text:     i.Text,
						Time:     time.Unix(i.Timestamp, 0)},
					)
				}
			}
		}
	}
	for u, list := range result {
		s.l.Printf("username: %v", u)
		for i, m := range list {
			s.l.Printf("Message[%d] %v", i, m)
		}
	}
	//log.Printf("%#v", result)
	return result, nil
}

// recipient of dialog
func (s *service) recipient(users []goinsta.User) string {
	for _, u := range users {
		_users[u.ID] = u.Username
		if u.Username == s.api.Account.Username {
			continue
		}
		return u.Username
	}
	return ""
}
