package main

import "github.com/charlesbao/frpc"

func main() {
	err := frpc.RunClient("config.ini")
	panic(err)
}
