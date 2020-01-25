package main

import "github.com/charlesbao/frpc"

func main() {
	err := frpc.Run("config.ini")
	panic(err)
}
