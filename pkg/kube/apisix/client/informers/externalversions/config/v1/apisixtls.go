/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	configv1 "github.com/api7/ingress-controller/pkg/kube/apisix/apis/config/v1"
	versioned "github.com/api7/ingress-controller/pkg/kube/apisix/client/clientset/versioned"
	internalinterfaces "github.com/api7/ingress-controller/pkg/kube/apisix/client/informers/externalversions/internalinterfaces"
	v1 "github.com/api7/ingress-controller/pkg/kube/apisix/client/listers/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ApisixTLSInformer provides access to a shared informer and lister for
// ApisixTLSs.
type ApisixTLSInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ApisixTLSLister
}

type apisixTLSInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewApisixTLSInformer constructs a new informer for ApisixTLS type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewApisixTLSInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredApisixTLSInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredApisixTLSInformer constructs a new informer for ApisixTLS type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredApisixTLSInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApisixV1().ApisixTLSs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApisixV1().ApisixTLSs(namespace).Watch(context.TODO(), options)
			},
		},
		&configv1.ApisixTLS{},
		resyncPeriod,
		indexers,
	)
}

func (f *apisixTLSInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredApisixTLSInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *apisixTLSInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configv1.ApisixTLS{}, f.defaultInformer)
}

func (f *apisixTLSInformer) Lister() v1.ApisixTLSLister {
	return v1.NewApisixTLSLister(f.Informer().GetIndexer())
}
