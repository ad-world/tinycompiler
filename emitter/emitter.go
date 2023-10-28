package emitter

import "os"

type Emitter struct {
	FullPath string
	Header   string
	Code     string
}

func NewEmitter(fullPath string) *Emitter {
	e := &Emitter{
		FullPath: fullPath,
		Code:     "",
		Header:   "",
	}

	return e
}

func Emit(e *Emitter, code string) {
	e.Code += code
}

func EmitLine(e *Emitter, code string) {
	e.Code += code + "\n"
}

func HeaderLine(e *Emitter, code string) {
	e.Header += code + "\n"
}

func WriteFile(e *Emitter) {
	f, err := os.Create(e.FullPath)

	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(e.Header + e.Code)

	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}
