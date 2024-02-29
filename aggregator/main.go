package main

import "log"

const httpServerAddr = ":3001"

func main() {
	log.Fatal(HTTPServer(httpServerAddr))
}
