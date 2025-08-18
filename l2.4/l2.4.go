package main

func main() {
	ch := make(chan int)
	go func() {
		// defer close(ch) можно вот так решить проблему
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	for n := range ch {
		println(n) //1. = 0, 2. = 1,,,,,, 10. = 9
	}
	// DEADLOCK потому что пытается читать с канала хотя уже не пишет в него
}
