package clipboard

import (
	"context"

	"golang.design/x/clipboard"
)

type Clipboard struct {
	Data string
}

func (c *Clipboard) clipboardInit() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}

func (c *Clipboard) ReadClipboard() string{
	c.clipboardInit()
	return string(clipboard.Read(clipboard.FmtText))
}

func (c *Clipboard) WriteClipboard(text string) {
	c.clipboardInit();
	clipboard.Write(clipboard.FmtText, []byte(text));
}

func (c *Clipboard) WatchClipboard() string{
	c.clipboardInit();

	ch := clipboard.Watch(context.TODO(), clipboard.FmtText);
	for data := range ch {
		// print out clipboard data whenever it is changed
		return string(data)
	}
	return "error"
}