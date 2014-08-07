package fwk

import (
	"reflect"
)

type achan chan interface{}

type datastore struct {
	SvcBase
	store map[string]achan
}

func (ds *datastore) Configure(ctx Context) error {
	return nil
}

func (ds *datastore) Get(k string) (interface{}, error) {
	//fmt.Printf(">>> get(%v)...\n", k)
	ch, ok := ds.store[k]
	if !ok {
		return nil, Errorf("Store.Get: no such key [%v]", k)
	}
	v := <-ch
	ch <- v
	//fmt.Printf("<<< get(%v, %v)...\n", k, v)
	return v, nil
}

func (ds *datastore) Put(k string, v interface{}) error {
	//fmt.Printf(">>> put(%v, %v)...\n", k, v)
	ds.store[k] <- v
	//fmt.Printf("<<< put(%v, %v)...\n", k, v)
	return nil
}

func (ds *datastore) Has(k string) bool {
	_, ok := ds.store[k]
	return ok
}

func (ds *datastore) StartSvc(ctx Context) error {
	ds.store = make(map[string]achan)
	return nil
}

func (ds *datastore) StopSvc(ctx Context) error {
	ds.store = nil
	return nil
}

func init() {
	Register(reflect.TypeOf(datastore{}),
		func(typ, name string, mgr App) (Component, error) {
			return &datastore{
				SvcBase: NewSvc(typ, name, mgr),
				store:   make(map[string]achan),
			}, nil
		},
	)
}

// interface tests
var _ Store = (*datastore)(nil)

// EOF
