/*
Copyright 2022.

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

package controllers

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	podchaosv1alpha1 "github.com/perithompson/podchaosmonkey/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeClient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetInterval(t *testing.T) {
	type args struct {
		interval string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "5m",
			args: args{
				interval: "5m",
			},
			want:    time.Duration(5 * time.Minute),
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				interval: "",
			},
			want:    time.Duration(30 * time.Second),
			wantErr: false,
		},
		{
			name: "15s",
			args: args{
				interval: "15s",
			},
			want:    time.Duration(15 * time.Second),
			wantErr: false,
		},
	}
	g := NewWithT(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMinInterval(tt.args.interval)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(got).Should(Equal(tt.want))
			}
		})
	}
}

func TestMonkeyReconciler_Reconcile(t *testing.T) {
	clientBuilder := fakeClient.NewClientBuilder()
	clientBuilder.WithObjects(&podchaosv1alpha1.Monkey{})

	fakeScheme, err := podchaosv1alpha1.SchemeBuilder.Build()
	if err != nil {
		Panic()
	}
	c := clientBuilder.WithScheme(fakeScheme).Build()
	g := NewWithT(t)
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx context.Context
		req ctrl.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ctrl.Result
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Client: c,
				Scheme: fakeScheme,
			},
			args:    args{},
			want:    ctrl.Result{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MonkeyReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			got, err := r.Reconcile(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(got).Should(Equal(tt.want))
			}
		})
	}
}

func TestMonkeyReconciler_UpdateStatus(t *testing.T) {
	clientBuilder := fakeClient.NewClientBuilder()
	clientBuilder.WithObjects(&podchaosv1alpha1.Monkey{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "workloads",
		},
		Spec: podchaosv1alpha1.MonkeySpec{
			Noop:      false,
			Interval:  "5m",
			Namespace: "workloads",
			Selector:  metav1.LabelSelector{},
		},
		Status: podchaosv1alpha1.MonkeyStatus{},
	})

	fakeScheme, err := podchaosv1alpha1.SchemeBuilder.Build()
	if err != nil {
		Panic()
	}
	c := clientBuilder.WithScheme(fakeScheme).Build()
	g := NewWithT(t)
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx    context.Context
		monkey *podchaosv1alpha1.Monkey
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ctrl.Result
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Client: c,
				Scheme: fakeScheme,
			},
			args: args{
				ctx: context.Background(),
				monkey: &podchaosv1alpha1.Monkey{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "workloads",
					},
					Spec: podchaosv1alpha1.MonkeySpec{
						Noop:      false,
						Interval:  "5m",
						Namespace: "workloads",
						Selector:  metav1.LabelSelector{},
					},
					Status: podchaosv1alpha1.MonkeyStatus{},
				},
			},
			want: ctrl.Result{
				RequeueAfter: time.Duration(5 * time.Minute),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MonkeyReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			got, err := r.UpdateStatus(tt.args.ctx, tt.args.monkey)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(got).Should(Equal(tt.want))
			}
		})
	}
}
