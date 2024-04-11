package assets

import (
	"log"

	"github.com/webability-go/xmodules/base"
)

const (
	MODULEID = "translation"
	VERSION  = "0.0.1"
)

type ModuleEntries struct {
	TranslateOne func(input string) (string, error)
}

func GetEntries(logger *log.Logger) *ModuleEntries {
	me := base.GetEntries(logger, MODULEID)
	if me == nil {
		return nil
	}
	lme, ok := me.(*ModuleEntries)
	if !ok {
		return nil
	}
	return lme
}
