package inbox

import (
	"fmt"

	"github.com/lab5e/spancli/pkg/commonopt"
)

type listInboxCmd struct {
	ID   commonopt.CollectionAndDevice
	List commonopt.ListOptions
}

func (*listInboxCmd) Execute([]string) error {
	fmt.Println("List inbox")
	return nil
}
