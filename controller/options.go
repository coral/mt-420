package controller

type Option interface {
	Run(c *Controller, events <-chan string, end chan bool) string
	Name() string
}
