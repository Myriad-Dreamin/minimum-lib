package rbac

import (
	_ "fmt"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"sync"
	"time"
)

var initalLifeTime = 10
var tickDuration = time.Second

var maxDescriptorCount int = 15
var allocDescriptorCount int

func SetInitalLifeTime(i int) {
	initalLifeTime = i
}

func SetTickDuration(t time.Duration) {
	tickDuration = t
}

func SetMaxDescriptorCount(c int) {
	maxDescriptorCount = c
}

var fileMutex sync.Mutex

type EnforceDescriptor struct {
	lifetime  int
	reference int
	*casbin.SyncedEnforcer
}

func newDescriptor(e *casbin.SyncedEnforcer) *EnforceDescriptor {
	return &EnforceDescriptor{
		lifetime:       initalLifeTime,
		reference:      1,
		SyncedEnforcer: e,
	}
}

var cachedEnforcers = make(map[string]*EnforceDescriptor)

func Accquire(path string) (*casbin.SyncedEnforcer, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	if e, ok := cachedEnforcers[path]; ok {
		e.lifetime += 2
		if e.lifetime > initalLifeTime {
			e.lifetime = initalLifeTime
		}
		e.reference++
		return e.SyncedEnforcer, nil
	}

	if maxDescriptorCount <= allocDescriptorCount {
		// is it ok ?
		return casbin.NewSyncedEnforcer(rbacModel, fileadapter.NewAdapter(path))
	} else {
		allocDescriptorCount++
		enforcer, err := casbin.NewSyncedEnforcer(rbacModel, fileadapter.NewAdapter(path))
		if err != nil {
			return nil, err
		}
		e := newDescriptor(enforcer)
		cachedEnforcers[path] = e
		// fmt.Println("alocation...", e)
		go func() {
			ticker := time.NewTicker(tickDuration)
			var e *EnforceDescriptor
			for _ = range ticker.C {
				fileMutex.Lock()
				e = cachedEnforcers[path]
				e.lifetime--
				//fmt.Println(e)
				if e.reference <= 0 && e.lifetime <= 0 {
					allocDescriptorCount--
					delete(cachedEnforcers, path)
					fileMutex.Unlock()
					//fmt.Println("released", e)
					ticker.Stop()
					break
				}
				fileMutex.Unlock()
			}
		}()

		return e.SyncedEnforcer, nil
	}
}

func Release(path string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	if e, ok := cachedEnforcers[path]; ok {
		e.lifetime--
		e.reference--
	}
	return
}
