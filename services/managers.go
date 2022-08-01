package services

import (
	"github.com/fiure-cms/fiure-web/internal/managers"
	"github.com/uretgec/go-sonic/sonic"
)

var Bbm *managers.PostManager
var Pm *managers.PageManager
var Sm *managers.SearchManager

func SetupManagers() {
	Bbm = managers.NewPostManager(Sr, false)
	Pm = managers.NewPageManager(Sr, false)
	Sm = managers.NewSearchManager(Ss.Clients[sonic.ChannelSearch])
}
