package utils

type UUIDGenerator struct {
	idGen        uint32
	internalChan chan uint32
}

func NewUUIDGenerator(max uint32, cacheLen uint32) *UUIDGenerator {
	gen := &UUIDGenerator{
		idGen:        0,
		internalChan: make(chan uint32, cacheLen),
	}
	gen.startGen(max)
	return gen
}

//开启 goroutine, 把生成的数字形式的UUID放入缓冲管道
func (g *UUIDGenerator) startGen(max uint32) {
	go func() {
		for {
			if g.idGen == max {
				g.idGen = 1
			} else {
				g.idGen += 1
			}
			g.internalChan <- g.idGen
		}
	}()
}

// 获取uint32形式的UUID
func (g *UUIDGenerator) GetUint32() uint32 {
	return <-g.internalChan
}
