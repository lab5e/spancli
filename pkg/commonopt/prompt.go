package commonopt

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
)

type NoPrompt struct {
	YesIAmSure bool `long:"yes-i-am-sure" description:"disable prompt for 'are you sure'"`
}

func (p *NoPrompt) Check() bool {
	if !p.YesIAmSure {
		rand.Seed(time.Now().UnixNano())
		verify := fmt.Sprintf("yes %04d", rand.Intn(9999))

		fmt.Printf("\n%s\n\n", text.Colors{text.BgRed, text.FgWhite}.Sprint("*** D A N G E R ***"))
		fmt.Printf("enter '%s' to confirm: ", verify)
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		return text == (verify + "\n")
	}
	return true
}
