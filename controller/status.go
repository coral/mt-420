package controller

import "golang.org/x/exp/errors/fmt"

type Status struct {
}

func (m *Status) Run(c *Controller, events <-chan string, end chan bool) string {
	for {
		select {
		case <-end:
			return "status"
		case <-events:
			fmt.Println(events)
		default:

		}
	}
	return "status"

}

func (m *Status) Name() string {
	return "status"
}
