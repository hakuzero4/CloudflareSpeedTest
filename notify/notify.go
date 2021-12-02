package notify

import (
	"fmt"
	"strings"
)

type Notify interface {
	SendMsg(string) error
}

var notifyServers []Notify

func AddServer(n Notify) {
	notifyServers = append(notifyServers, n)
}

func SendMsg(m string) {
	var errstrings []string
	for _, sv := range notifyServers {
		if err := sv.SendMsg(m); err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}
	fmt.Println(fmt.Errorf(strings.Join(errstrings, "\n")))
}
