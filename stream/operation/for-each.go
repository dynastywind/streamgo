package operation

func ForEachParallel(arr []interface{}, num int, consumer func(interface{})) {
	length := len(arr)
	ch := make(chan int, length)
	for i, item := range arr {
		go func(data interface{}) {
			consumer(data)
			ch <- 1
		}(item)
		if (i+1)%num == 0 || i == length-1 {
			for j := 0; j <= i%num; j++ {
				<-ch
			}
		}
	}
}
