package web

import (
	"io"
)

func Pack() {

}

func UnPack(reader io.Reader) {
	io.ReadAll(reader)
}