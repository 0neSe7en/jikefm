package jike
import (
	"github.com/mdp/qrterminal"
	"os"
)

func GenerateQRCode(uuid string) {
	url := `jike://page.jk/web?url=https%3A%2F%2Fruguoapp.com%2Faccount%2Fscan%3Fuuid%3D` + uuid + `&displayHeader=false&displayFooter=false`
	config := qrterminal.Config{
		Level: qrterminal.L,
		Writer: os.Stdout,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 2,
	}
	qrterminal.GenerateWithConfig(url, config)
}
