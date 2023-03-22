package dump

import "github.com/davecgh/go-spew/spew"

func init() {
	spew.Config = spew.ConfigState{
		Indent: "\t",
	}
}
