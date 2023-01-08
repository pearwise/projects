package peer

import (
	"syscall/js"

)

func Monitor(data js.Value) {
	data.Index(0).String()
}