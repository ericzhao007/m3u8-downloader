package spiders

func WithMemoryStoreEngine() WithFunc {
	return func(s *Spiders) {
		s.tsData = NewMemStoreEngine()
	}
}

func WithDiskStoreEngine() WithFunc {
	return func(s *Spiders) {
		s.tsData = NewDiskStoreEngine("./work")
	}
}

func WithWorkerNum(workerNum int) WithFunc {
	return func(s *Spiders) {
		s.workerNum = workerNum
	}
}
