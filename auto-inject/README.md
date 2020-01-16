
## auto-inject

#### Usage

```go
package main

import "github.com/Myriad-Dreamin/minimum-lib/auto-inject"
import "fmt"

type A int
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

func main() {

	s, err := auto_inject.ParseFlatSource(MInject{1, 2, 3, 4, 5})
	if err != nil {
		panic(err)
	}

	var a InjectTargetA
	l, err := s.Bind(&a)
	if err != nil {
		panic(err)
	}
	fmt.Println(a, l)
	var b InjectTargetB
	l, err = s.Bind(&b)
	if err != nil {
		panic(err)
	}
	fmt.Println(b, l)
}
```

```plain
GOROOT=/home/kamiyoru/work/go #gosetup
GOPATH=/home/kamiyoru/go #gosetup
/home/kamiyoru/work/go/bin/go test -c -o /tmp/___go_test_github_com_Myriad_Dreamin_minimum_lib_auto_inject -gcflags "all=-N -l" github.com/Myriad-Dreamin/minimum-lib/auto-inject #gosetup
/home/kamiyoru/work/go/bin/go tool test2json -t /opt/GoLand-2019.2.2/plugins/go/lib/dlv/linux/dlv --listen=localhost:38329 --headless=true --api-version=2 exec /tmp/___go_test_github_com_Myriad_Dreamin_minimum_lib_auto_inject -- -test.v #gosetup
API server listening at: 127.0.0.1:38329
=== RUN   TestInjectFlat
auto_inject.A
    A2 {{A2  auto_inject.A injector:"a2" 16 [2] false} <auto_inject.A Value>}
     {{AA  auto_inject.A  0 [0] false} <auto_inject.A Value>}
    AA {{AA  auto_inject.A  0 [0] false} <auto_inject.A Value>}
    a {{A  auto_inject.A injector:"a" 8 [1] false} <auto_inject.A Value>}
    A {{A  auto_inject.A injector:"a" 8 [1] false} <auto_inject.A Value>}
    a2 {{A2  auto_inject.A injector:"a2" 16 [2] false} <auto_inject.A Value>}
auto_inject.B
     {{B  auto_inject.B  24 [3] true} <auto_inject.B Value>}
    B {{B  auto_inject.B  24 [3] true} <auto_inject.B Value>}
auto_inject.C
    c {{C  auto_inject.C injector:"c" 32 [4] true} <auto_inject.C Value>}
    C {{C  auto_inject.C injector:"c" 32 [4] true} <auto_inject.C Value>}
     {{C  auto_inject.C injector:"c" 32 [4] true} <auto_inject.C Value>}
{1 2 2} []
{2 3 3 1 0} [A5]
--- PASS: TestInjectFlat (0.00s)
PASS

Debugger finished with exit code 0

```

