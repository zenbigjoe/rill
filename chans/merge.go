package chans

import (
	"fmt"
	"sync"
)

func fastMerge[A any](ins []<-chan A) <-chan A {
	// len(ins) must be between 2 and 5

	remaining := len(ins)
	for len(ins) < 5 {
		ins = append(ins, nil)
	}

	out := make(chan A)

	go func() {
		defer close(out)

		var a A
		var ok bool
		var i int

		for {
			if remaining == 0 {
				return
			}

			select {
			case a, ok = <-ins[0]:
				i = 0
			case a, ok = <-ins[1]:
				i = 1
			case a, ok = <-ins[2]:
				i = 2
			case a, ok = <-ins[3]:
				i = 3
			case a, ok = <-ins[4]:
				i = 4
			}

			if !ok {
				remaining--
				ins[i] = nil
				continue
			}

			out <- a
		}
	}()

	return out
}

func slowMerge[A any](ins []<-chan A) <-chan A {
	out := make(chan A)

	var wg sync.WaitGroup
	for _, in := range ins {
		in1 := in
		wg.Add(1)
		go func() {
			defer wg.Done()
			for x := range in1 {
				out <- x
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Merge[A any](ins ...<-chan A) <-chan A {
	switch len(ins) {
	case 0:
		return nil
	case 1:
		return ins[0]
	case 2, 3, 4, 5:
		return fastMerge(ins)
	default:
		return slowMerge(ins)
	}
}

// todo: not sure if this is a good idea
// todo: introduce concurrency here?
func Split2[A any](in <-chan A, f func(A) bool) (outTrue <-chan A, outFalse <-chan A) {
	if in == nil {
		return nil, nil
	}

	outT := make(chan A)
	outF := make(chan A)

	go func() {
		defer close(outT)
		defer close(outF)

		for a := range in {
			fmt.Println("split2", a)

			if f(a) {
				outT <- a
			} else {
				outF <- a
			}
		}
	}()

	return outT, outF
}
