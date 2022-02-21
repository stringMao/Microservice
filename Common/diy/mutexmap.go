package diy

import "sync"

//多线程安全的自定义map
type MutexMap struct {
	Data map[interface{}]interface{}
	lock *sync.RWMutex
}

func (m *MutexMap) Get(key interface{}) (value interface{}, ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v, o := m.Data[key]
	return v, o
}

func (m *MutexMap) Add(key, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.Data[key] = value
}

func (m *MutexMap) Delete(key interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	//if _, ok := m.Data[key] && ok{
	delete(m.Data, key)
	//}
}

func (m *MutexMap) Range(f func(key, value interface{}) bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for k, v := range m.Data {
		if f(k, v) {
			break
		}
	}
}
