package storage

import (
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
)

type StrategyProvider interface {
	SetStrategy(Strategy)
	GetStrategy() Strategy
}

type Storage interface {
	rest.StandardStorage
	rest.Scoper
	rest.Storage
	rest.ShortNamesProvider
	rest.TableConvertor
	registry.GenericStore
	StrategyProvider
}

type PredicateFunc = func(label labels.Selector, field fields.Selector) storage.SelectionPredicate
type NewObjectFunc = func() runtime.Object
type NewStorageFunc = func(
	gr schema.GroupResource,
	strategy Strategy,
	optsGetter generic.RESTOptionsGetter,
	tableConvertor rest.TableConvertor,
	newFunc, newListFunc NewObjectFunc,
) (Storage, error)
