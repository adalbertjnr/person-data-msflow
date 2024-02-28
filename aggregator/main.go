package main

import "log"

const httpServerAddr = ":3000"

func main() {
	log.Fatal(HTTPServer(httpServerAddr))
}
