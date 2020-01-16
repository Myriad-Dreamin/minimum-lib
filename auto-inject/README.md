
## auto-inject

```plain
=== RUN   TestInjectFlat
? {A  auto_inject.A injector:"a" 8 [1] false}
? {A2  auto_inject.A injector:"a2" 16 [2] false}
? {C  auto_inject.C injector:"c" 32 [4] true}
auto_inject.B
     {{B  auto_inject.B  24 [3] true} <auto_inject.B Value>}
    B {{B  auto_inject.B  24 [3] true} <auto_inject.B Value>}
auto_inject.C
    c {{C  auto_inject.C injector:"c" 32 [4] true} <auto_inject.C Value>}
    C {{C  auto_inject.C injector:"c" 32 [4] true} <auto_inject.C Value>}
     {{C  auto_inject.C injector:"c" 32 [4] true} <auto_inject.C Value>}
auto_inject.A
    a2 {{A2  auto_inject.A injector:"a2" 16 [2] false} <auto_inject.A Value>}
    A2 {{A2  auto_inject.A injector:"a2" 16 [2] false} <auto_inject.A Value>}
     {{AA  auto_inject.A  0 [0] false} <auto_inject.A Value>}
    AA {{AA  auto_inject.A  0 [0] false} <auto_inject.A Value>}
    a {{A  auto_inject.A injector:"a" 8 [1] false} <auto_inject.A Value>}
    A {{A  auto_inject.A injector:"a" 8 [1] false} <auto_inject.A Value>}
{1 2 2} []
{2 3 3 1 0} [A5]
--- PASS: TestInjectFlat (0.00s)
PASS
ok  	_/home/kamiyoru/work/gosrc/src/github.com/Myriad-Dreamin/auto-inject	0.003s

Process finished with exit code 0
```

