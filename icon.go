package bui

import (
	"crypto/md5"
	"fmt"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
	"io/ioutil"
	"os"
	"path/filepath"
)

func IconFromBytes(iconBytes []byte) (win.HICON, error) {
	md5Result := md5.Sum(iconBytes)
	filename := fmt.Sprintf("%x.ico", md5Result)
	iconPath := filepath.Join(os.TempDir(), TempDir, filename)
	err := ioutil.WriteFile(iconPath, iconBytes, 0644)
	if err != nil {
		return win.HICON(0), err
	}

	icon := win.LoadImage(
		0,
		windows.StringToUTF16Ptr(iconPath),
		win.IMAGE_ICON,
		0,
		0,
		win.LR_DEFAULTSIZE|win.LR_LOADFROMFILE)

	return win.HICON(icon), nil
}
