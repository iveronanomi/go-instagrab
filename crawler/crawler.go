package crawler

import (
	"log"

	"github.com/ahmdrz/goinsta"
)

type service struct {
	api *goinsta.Instagram
	l   *log.Logger
}

// New crawler
func New(api *goinsta.Instagram, l *log.Logger) *service {
	return &service{
		api: api,
		l:   l,
	}
}
