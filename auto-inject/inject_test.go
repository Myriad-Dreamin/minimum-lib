package auto_inject

import (
	"fmt"
	"testing"
)

type A interface {}
type B int
type C int

type MInject struct {
	AA A
	A  A `injector:"a"`
	A2 A `injector:"a2"`
	B
	C `injector:"c"`
}

type InjectTargetA struct {
	AA A
	A  A `injector:"a"`
	A2 A `injector:"a"`
}

type InjectTargetB struct {
	AA A `injector:"a"`
	A  A `injector:"a2"`
	A2 A
	A4 A
	A5 A `injector:"a3"`
}

func TestInjectFlat(t *testing.T) {
	s, err := ParseFlatSource(MInject{1, 2, 3, 4, 5})
	if err != nil {
		t.Fatal(err)
	}
	for t, x := range s {
		fmt.Println(t)
		for d, y := range x {
			fmt.Println("   ", d, y)
		}
	}
	var a InjectTargetA
	l, err := s.Bind(&a)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(a, l)
	var b InjectTargetB
	l, err = s.Bind(&b)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b, l)
}
