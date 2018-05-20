package gitmono

import (
	"github.com/ByteFlinger/tecos/backend"
)

func (d *GitMono) ListModules() []backend.ModuleData {

	return []backend.ModuleData{}

}

func (d *GitMono) Cleanup() {
	close(d.quitChan)
	<-d.doneChan
}
