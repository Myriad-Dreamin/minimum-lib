package rbac

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAddPolicy(t *testing.T) {

	e, err := Accquire("./test.csv")
	if err != nil {
		t.Fatal(err)
	}
	added, err := e.AddPolicy("data2_admin", "data2", "write")
	if added == false || err != nil {
		t.Fatal(added, err)
	}
	if err = e.SavePolicy(); err != nil {
		t.Fatal(err)
	}
	Release("./test.csv")
}

func TestDeletePolicy(t *testing.T) {

	e, err := Accquire("./test.csv")
	if err != nil {
		t.Fatal(err)
	}
	removed, err := e.RemovePolicy("data2_admin", "data2", "write")
	if removed == false || err != nil {
		t.Fatal(removed, err)
	}
	if err := e.SavePolicy(); err != nil {
		t.Error(err)
	}
	Release("./test.csv")
}

func TestAccquireRelease(t *testing.T) {
	dur := time.Millisecond * 100
	SetTickDuration(dur)
	e, err := Accquire("./test.csv")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(e.GetPolicy())
	time.Sleep(dur * 5)
	e2, err := Accquire("./test.csv")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%p %p\n", e, e2)
	Release("./test.csv")
	time.Sleep(dur * 8)
	Release("./test.csv")
	time.Sleep(dur * 1)
	SetTickDuration(time.Second)
}

func BenchmarkAccquire(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Accquire("./test.csv")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAccquireReleaseMultiThread(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 1000; j++ {
			wg.Add(1)
			go func() {
				e, err := Accquire("./test.csv")
				if err != nil {
					b.Fatal(err)
				}

				added, err := e.AddPolicy("data2_admin", "data2", "write")
				if added == false || err != nil {
					b.Fatal(added, err)
				}
				if err := e.SavePolicy(); err != nil {
					b.Error(err)
				}
				removed, err := e.RemovePolicy("data2_admin", "data2", "write")
				if removed == false || err != nil {
					b.Fatal(removed, err)
				}
				if err := e.SavePolicy(); err != nil {
					b.Error(err)
				}
				Release("./test.csv")
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkAccquireRelease(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e, err := Accquire("./test.csv")
		if err != nil {
			b.Fatal(err)
		}

		added, err := e.AddPolicy("data2_admin", "data2", "write")
		if added == false || err != nil {
			b.Fatal(added, err)
		}
		if err := e.SavePolicy(); err != nil {
			b.Error(err)
		}
		removed, err := e.RemovePolicy("data2_admin", "data2", "write")
		if removed == false || err != nil {
			b.Fatal(removed, err)
		}
		if err := e.SavePolicy(); err != nil {
			b.Error(err)
		}
		Release("./test.csv")
	}
}
