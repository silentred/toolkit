package rotator

import (
	"log"
	"testing"
)

func Test_FileSizeRotator(t *testing.T) {
	spliter := NewSizeSpliter(defaultSize)
	fileRotator := NewFileRotator("", "app", "log", spliter)
	l := log.New(fileRotator, "test", log.LstdFlags)

	for i := 0; i < 30; i++ {
		l.Printf("tess %d \n", i)
	}
}

func Test_FileDayRotator(t *testing.T) {
	spliter := NewDaySpliter()
	fileRotator := NewFileRotator("", "app", "log", spliter)
	l := log.New(fileRotator, "test", log.LstdFlags)

	for i := 0; i < 30; i++ {
		l.Printf("tess %d \n", i)
	}
}
