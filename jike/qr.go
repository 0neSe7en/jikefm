package jike

import (
	"os"

	"github.com/mdp/qrterminal"
)

func GenerateQRCode(uuid string) {
	url := `jike://page.jk/web?url=https%3A%2F%2Fruguoapp.com%2Faccount%2Fscan%3Fuuid%3D` + uuid + `&displayHeader=false&displayFooter=false`
	config := qrterminal.Config{
		Level:          qrterminal.L,
		Writer:         os.Stdout,
		HalfBlocks:     true,
		BlackChar:      qrterminal.BLACK_BLACK,
		WhiteBlackChar: qrterminal.WHITE_BLACK,
		WhiteChar:      qrterminal.WHITE_WHITE,
		BlackWhiteChar: qrterminal.BLACK_WHITE,
		QuietZone:      2,
	}
	qrterminal.GenerateWithConfig(url, config)
}
