package main

func main() {
	s := NewServer("localhost", 8888)
	s.Start()
}
