package module

import (
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/logger"
	"os"
	"path"
	"reflect"
)

type Resource struct {
	Name  string
	Value interface{}
	Proto reflect.Type
	Type  reflect.Type
}

type Moduler interface {
	GetNamespace() string
	Requires() (r []string)
	Provides() (r []Resource)
	AfterInstall(m Module) error
}

type PModuler interface {
	Moduler
	Provide(Name string, Value interface{}) error
}

type BaseModuler struct {
	Resources []Resource
	Namespace string
}

func (b *BaseModuler) GetNamespace() string {
	return b.Namespace
}

func (b *BaseModuler) Requires() (r []string) {
	return nil
}

func (b *BaseModuler) Provides() []Resource {
	return b.Resources
}

func (b *BaseModuler) Provide(Name string, Value interface{}) (err error) {
	t := reflect.TypeOf(Value)
	for _, res := range b.Resources {
		if res.Type == t {
			if len(res.Name) == 0 {
				return fmt.Errorf("shadowing anonymous resource namespace %v", Name)
			}
			if res.Name == Name {
				return fmt.Errorf("duplicate resource namespace %v", Name)
			}
		}
	}
	b.Resources = append(b.Resources, Resource{Name: Name, Value: Value, Type: t})
	return nil
}

func (b *BaseModuler) ProvideNamedImpl(Name string, protoPtr, impl interface{}) (err error) {
	proto := reflect.TypeOf(protoPtr).Elem()
	if err = AssertHasImpl(impl, proto); err != nil {
		return err
	}
	t := reflect.TypeOf(impl)
	for _, res := range b.Resources {
		if res.Type == t {
			if len(res.Name) == 0 {
				return fmt.Errorf("shadowing anonymous resource namespace %v", Name)
			}
			if res.Name == Name {
				return fmt.Errorf("duplicate resource namespace %v", Name)
			}
		}
	}
	b.Resources = append(b.Resources, Resource{Name: Name, Value: impl, Proto: proto, Type: t})
	return nil
}

func (b *BaseModuler) ProvideImpl(protoPtr, impl interface{}) (err error) {
	proto := reflect.TypeOf(protoPtr).Elem()
	if err = AssertHasImpl(impl, proto); err != nil {
		return err
	}
	t := reflect.TypeOf(impl)
	for _, res := range b.Resources {
		if res.Type == t {
			return fmt.Errorf("duplicate resource namespace %v", res.Name)
		}
	}
	b.Resources = append(b.Resources, Resource{Name: "", Value: impl, Proto: proto, Type: t})
	return nil
}

func (b *BaseModuler) Replace(Name string, Value interface{}) (err error) {
	for i, res := range b.Resources {
		if res.Name == Name {
			b.Resources[i].Value = Value
			b.Resources[i].Type = reflect.TypeOf(Value)
			return nil
		}
	}
	b.Resources = append(b.Resources, Resource{Name: Name, Value: Value, Type: reflect.TypeOf(Value)})
	return nil
}

func (b BaseModuler) AfterInstall(m Module) error {
	return nil
}

type Module map[string]map[reflect.Type]interface{}

func (m Module) Require(namespace string) interface{} {
	return m.requireNamed(namespace)
}

func (m Module) RequireNamedImpl(namespace string, protoPtr interface{}) interface{} {
	return m.requireNamedImpl(namespace, reflect.TypeOf(protoPtr).Elem())
}

func (m Module) RequireImpl(protoPtr interface{}) interface{} {
	return m.requireNamedImpl("", reflect.TypeOf(protoPtr).Elem())
}

func (m Module) requireNamed(namespace string) interface{} {
	if m == nil {
		return nil
	}
	sm := m[namespace]
	if sm == nil {
		return nil
	}
	return sm[nil]
}

func (m Module) requireNamedImpl(namespace string, proto reflect.Type) interface{} {
	if m == nil {
		return nil
	}
	sm := m[namespace]
	if sm == nil {
		return nil
	}
	if proto.Kind() == reflect.Interface {
		for t, v := range sm {
			if t.Implements(proto) {
				return v
			}
		}
		return nil
	} else {
		return sm[proto]
	}

}

func (m Module) ProvideImpl(protoPtr, impl interface{}) (err error) {
	if err = AssertHasImpl(impl, reflect.TypeOf(protoPtr).Elem()); err != nil {
		return err
	}
	return m.ProvideWithCheck("", impl)
}

func (m Module) ProvideNamedImpl(namespace string, protoPtr, impl interface{}) (err error) {
	proto := reflect.TypeOf(protoPtr).Elem()
	if err = AssertHasImpl(impl, proto); err != nil {
		return err
	}
	_, ov, err := m.provideWithCheck(namespace, proto, impl)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "warning: overriding Type: %v\n", err)
	}
	if ov != nil {
		return fmt.Errorf("install provide error: conflict on namespace:%v, impl:%T", namespace, impl)
	}
	return nil
}

func (m Module) ProvideNamedImplT(namespace string, proto reflect.Type, impl interface{}) (err error) {
	if err = AssertHasImpl(impl, proto); err != nil {
		return err
	}
	_, ov, err := m.provideWithCheck(namespace, proto, impl)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "warning: overriding Type: %v\n", err)
	}
	if ov != nil {
		return fmt.Errorf("install provide error: conflict on namespace:%v, impl:%T", namespace, impl)
	}
	return nil
}

func (m Module) Provide(namespace string, impl interface{}) (tov, ov interface{}) {
	tov, ov, _ = m.provideWithCheck(namespace, reflect.TypeOf(impl), impl)
	return
}

func (m Module) ProvideWithCheck(namespace string, impl interface{}) error {
	tov, _, err := m.provideWithCheck(namespace, reflect.TypeOf(impl), impl)
	if err != nil {
		return err
	}
	if tov != nil {
		return fmt.Errorf("install provide error: conflict on impl:%T", impl)
	}
	return nil
}

func (m Module) provideWithCheck(namespace string, t reflect.Type, impl interface{}) (tov, ov interface{}, err error) {
	if m == nil {
		panic(fmt.Errorf("nil module container"))
	}
	sm := m[""]
	if sm == nil {
		sm = make(map[reflect.Type]interface{})
		m[""] = sm
	}
	var ok bool
	tov, ok = sm[t]
	sm[t] = impl
	if ok {
		err = fmt.Errorf("already provided a impl of %T", impl)
	}
	if len(namespace) == 0 {
		return
	}
	sm = m[namespace]
	if sm == nil {
		sm = make(map[reflect.Type]interface{})
		m[namespace] = sm
	}
	switch len(sm) {
	case 0, 1:
		sm[nil] = impl
	case 2:
		ov, ok = sm[t]
		if ok {
			sm[nil] = impl
		} else {
			sm[nil] = nil
		}
	case 3:
		delete(sm, nil)
		fallthrough
	default:
		ov = sm[t]
	}
	sm[t] = impl
	return
}

func AssertHasImpl(impl interface{}, proto reflect.Type) error {
	if proto == nil {
		return nil
	}
	if proto.Kind() == reflect.Interface {
		if !reflect.TypeOf(impl).Implements(proto) {
			return fmt.Errorf("impl %T does not implement %v", impl, proto)
		}
		return nil
	}
	if proto != reflect.TypeOf(impl) {
		return fmt.Errorf("impl %T has no type %v", impl, proto)
	}
	return nil
}

func (m Module) Install(moduler Moduler) (err error) {
	namespace := moduler.GetNamespace()
	for _, res := range moduler.Provides() {
		err := m.ProvideNamedImplT(path.Join(namespace, res.Name), res.Proto, res.Value)
		if err != nil {
			return err
		}
	}

	for _, req := range moduler.Requires() {
		sm := m[req]
		if sm == nil || sm[nil] == nil {
			return fmt.Errorf("does not meets the requirement %v", req)
		}
	}

	err = moduler.AfterInstall(m)
	if err != nil {
		return
	}
	return nil
}

func (m Module) Debug(logger logger.Logger) {
	for k, v := range m {
		logger.Debug("module installed", "path", k, "namespace", reflect.TypeOf(v))
	}
}
