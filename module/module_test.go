package module

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"reflect"
	"testing"
)

type myModule struct {
	b *int
}

func (m myModule) GetNamespace() string {
	return ""
}

func (m myModule) Requires() (r []string) {
	return nil
}

func (m myModule) Provides() (r []Resource) {
	return []Resource{
		{Name: "mymodule.a", Value: 1, Type: reflect.TypeOf(1)},
		{Name: "mymodule.b", Value: m.b, Type: reflect.TypeOf(m.b)},
	}
}

func (m myModule) AfterInstall(Module) error {
	return nil
}

type subModule struct {
	b *int
}

func (m subModule) GetNamespace() string {
	return ""
}

func (m subModule) Requires() (r []string) {
	return []string{"mymodule.a", "mymodule.b"}
}

func (m subModule) Provides() (r []Resource) {
	return nil
}

func (m subModule) AfterInstall(mm Module) error {
	fmt.Println(mm.Require("mymodule.a"), *mm.Require("mymodule.b").(*int))
	return nil
}

func TestModule_Install(t *testing.T) {
	var m = make(Module)
	var mm = &myModule{b: new(int)}
	err := m.Install(mm)
	if err != nil {
		t.Fatal(err)
	}
	var mmm = &subModule{}
	err = m.Install(mmm)
	if err != nil {
		t.Fatal(err)
	}

	*mm.b = 0
	assert.Equal(t, 1, m.Require("mymodule.a"))
	assert.Equal(t, 0, *m.Require("mymodule.b").(*int))
	*mm.b = 1
	assert.Equal(t, 1, m.Require("mymodule.a"))
	assert.Equal(t, 1, *m.Require("mymodule.b").(*int))
}

func TestModule_DynamicProvideInstall(t *testing.T) {
	var m = make(Module)

	var mmm = new(BaseModuler)

	assert.NoError(t, mmm.Provide("global_var/1", 1))
	assert.Error(t, mmm.Provide("global_var/1", 2))

	err := m.Install(mmm)
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, mmm.ProvideImpl(new(float64), 3.))
	assert.Error(t, mmm.ProvideImpl(new(float64), 4.))
	assert.Error(t, mmm.ProvideImpl(new(int), 5))

	mmm.Namespace = "some-namespace"

	err = m.Install(mmm)
	if err != nil {
		t.Fatal(err)
	}

	//*mm.b = 0
	assert.Equal(t, 1, m.Require("global_var/1"))
	assert.Equal(t, 3., m.RequireImpl(new(float64)))
	assert.Equal(t, 1, m.Require("some-namespace/global_var/1"))
	assert.Equal(t, 3., m.RequireImpl(new(float64)))
	//assert.Equal(t, 0, *m.Require("mymodule.b").(*int))
	//*mm.b = 1
	//assert.Equal(t, 1, m.Require("mymodule.a"))
	//assert.Equal(t, 1, *m.Require("mymodule.b").(*int))
}

func TestModule_ProvideInner(t *testing.T) {
	var m = make(Module)
	assert.Equal(t, 0, len(m["1"]))

	var doOnKey = func(k string) {
		tov, ov := m.Provide(k, 1)
		assert.Equal(t, 1, m.Require(k))
		assert.Equal(t, 1, m.RequireImpl(new(int)))
		assert.Equal(t, 1, m.RequireNamedImpl(k, new(int)))
		//assert.Equal(t, nil, tov)
		assert.Equal(t, nil, ov)
		assert.Equal(t, 2, len(m[k]))

		// unsafe replacement
		tov, ov = m.Provide(k, 2)
		assert.Equal(t, 2, m.Require(k))
		assert.Equal(t, 2, m.RequireImpl(new(int)))
		assert.Equal(t, 2, m.RequireNamedImpl(k, new(int)))
		assert.Equal(t, 1, tov)
		assert.Equal(t, 1, ov)
		assert.Equal(t, 2, len(m[k]))

		// provide different type in the same namespace
		tov, ov = m.Provide(k, 2.)
		assert.Equal(t, nil, m.Require(k))
		assert.Equal(t, 2, m.RequireImpl(new(int)))
		assert.Equal(t, 2, m.RequireNamedImpl(k, new(int)))
		assert.Equal(t, 2., m.RequireNamedImpl(k, new(float64)))
		//assert.Equal(t, nil, tov)
		assert.Equal(t, nil, ov)
		assert.Equal(t, 3, len(m[k]))
	}

	doOnKey("1")
	doOnKey("2")
}

func TestModule_ProvideWithCheck(t *testing.T) {
	wr, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	r := io.TeeReader(wr, os.Stderr)
	os.Stderr = w

	var m = make(Module)

	assert.NoError(t, m.ProvideImpl(new(int), 1))
	assert.Equal(t, 1, m.RequireImpl(new(int)))
	assert.Error(t, m.ProvideImpl(new(int), 2))

	assert.Equal(t, 2, m.RequireImpl(new(int)))
	assert.NoError(t, m.ProvideNamedImpl("1", new(int), 3))
	assert.Equal(t, 3, m.RequireNamedImpl("1", new(int)))

	var b = make([]byte, 512)
	n, err := r.Read(b)
	assert.NotEqual(t, 0, n)
	assert.NoError(t, err)

	m = make(Module)
	assert.Equal(t, nil, m.RequireNamedImpl("1", new(int)))
	assert.Equal(t, nil, m.RequireNamedImpl("2", new(int)))
	assert.NoError(t, m.ProvideNamedImpl("1", new(int), 1))
	assert.Equal(t, 1, m.RequireImpl(new(int)))
	assert.Equal(t, 1, m.RequireNamedImpl("1", new(int)))
	assert.Equal(t, nil, m.RequireNamedImpl("2", new(int)))
	assert.NoError(t, m.ProvideNamedImpl("2", new(int), 2))

	n, err = r.Read(b)
	assert.NotEqual(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, 2, m.RequireImpl(new(int)))
	assert.Equal(t, 1, m.RequireNamedImpl("1", new(int)))
	assert.Equal(t, 2, m.RequireNamedImpl("2", new(int)))
	assert.Error(t, m.ProvideNamedImpl("1", new(int), 3))

	n, err = r.Read(b)
	assert.NotEqual(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, 3, m.RequireImpl(new(int)))
	assert.Equal(t, 3, m.RequireNamedImpl("1", new(int)))
	assert.Equal(t, 2, m.RequireNamedImpl("2", new(int)))
	assert.Error(t, m.ProvideNamedImpl("2", new(int), 4))

	n, err = r.Read(b)
	assert.NotEqual(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, 4, m.RequireImpl(new(int)))
	assert.Equal(t, 3, m.RequireNamedImpl("1", new(int)))
	assert.Equal(t, 4, m.RequireNamedImpl("2", new(int)))
}

type idleStringer struct {
}

func (i idleStringer) String() string {
	panic("implement me")
}

func TestModule_NotImplemented(t *testing.T) {

	m := make(Module)
	assert.Error(t, m.ProvideImpl(new(float64), 1))
	assert.Equal(t, nil, m.RequireImpl(new(float64)))
	assert.Error(t, m.ProvideNamedImpl("1", new(float64), 1))
	assert.Equal(t, nil, m.RequireImpl(new(float64)))
	assert.Equal(t, nil, m.RequireNamedImpl("1", new(float64)))
	assert.Error(t, m.ProvideImpl(new(fmt.Stringer), 1))
	assert.Equal(t, nil, m.RequireImpl(new(fmt.Stringer)))
	i := &idleStringer{}
	var _ fmt.Stringer = i
	assert.NoError(t, m.ProvideImpl(new(fmt.Stringer), i))
	assert.Equal(t, i, m.RequireImpl(new(fmt.Stringer)))
}
