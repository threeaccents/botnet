package wsync

//WaitGroup is
type WaitGroup struct {
	size int
	cur  int
}

//New is
func New(n int64) *WaitGroup {
	pool := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		channel <- struct{}{}
	}
	return &WaitGroup{
		size: n,
	}
}

//Add is
func (w *WaitGroup) Add(delta int) {

}

//Done is
func (w *WaitGroup) Done() {
	w.Add(-1)
}
