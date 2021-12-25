package main

func main() {
	r := InitRouter()
	panic(r.Run())
}
