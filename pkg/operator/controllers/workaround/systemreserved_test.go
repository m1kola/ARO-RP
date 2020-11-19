package workaround

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mcv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	fakemcoclient "github.com/openshift/machine-config-operator/pkg/generated/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ktesting "k8s.io/client-go/testing"

	"github.com/Azure/ARO-RP/pkg/util/cmp"
	utillog "github.com/Azure/ARO-RP/pkg/util/log"
	mock_dynamichelper "github.com/Azure/ARO-RP/pkg/util/mocks/dynamichelper"
	_ "github.com/Azure/ARO-RP/pkg/util/scheme"
)

func TestSystemreservedEnsure(t *testing.T) {
	tests := []struct {
		name                         string
		mcocli                       *fakemcoclient.Clientset
		mocker                       func(mdh *mock_dynamichelper.MockInterface)
		machineConfigPoolNeedsUpdate bool
		wantErr                      bool
	}{
		{
			name: "first time create",
			mcocli: fakemcoclient.NewSimpleClientset(&mcv1.MachineConfigPool{
				ObjectMeta: metav1.ObjectMeta{
					Name: "worker",
				},
			}),
			machineConfigPoolNeedsUpdate: true,
			mocker: func(mdh *mock_dynamichelper.MockInterface) {
				mdh.EXPECT().Ensure(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "nothing to be done",
			mcocli: fakemcoclient.NewSimpleClientset(&mcv1.MachineConfigPool{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "worker",
					Labels: map[string]string{labelName: labelValue},
				},
			}),
			mocker: func(mdh *mock_dynamichelper.MockInterface) {
				mdh.EXPECT().Ensure(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}
	err := mcv1.AddToScheme(scheme.Scheme)
	if err != nil {
		t.Error(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mdh := mock_dynamichelper.NewMockInterface(controller)
			sr := &systemreserved{
				mcocli: tt.mcocli,
				dh:     mdh,
				log:    utillog.GetLogger(),
			}

			var updated bool
			tt.mcocli.PrependReactor("update", "machineconfigpools", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
				updated = true
				return false, nil, nil
			})

			tt.mocker(mdh)
			if err := sr.Ensure(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("systemreserved.Ensure() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.machineConfigPoolNeedsUpdate != updated {
				t.Errorf("systemreserved.Ensure() updated %v, machineConfigPoolNeedsUpdate = %v", updated, tt.machineConfigPoolNeedsUpdate)
			}
		})
	}
}

func TestKubeletConfig(t *testing.T) {
	err := mcv1.AddToScheme(scheme.Scheme)
	if err != nil {
		t.Error(err)
	}
	sr := &systemreserved{}
	got, err := sr.kubeletConfig()
	if err != nil {
		t.Errorf("systemreserved.kubeletConfig() error = %v", err)
		return
	}
	want := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "machineconfiguration.openshift.io/v1",
			"kind":       "KubeletConfig",
			"metadata": map[string]interface{}{
				"creationTimestamp": nil,
				"labels": map[string]interface{}{
					"aro.openshift.io/limits": "",
				},
				"name": kubeletConfigName,
			},
			"spec": map[string]interface{}{
				"kubeletConfig": map[string]interface{}{
					"systemReserved": map[string]interface{}{
						"memory": "2000Mi",
					},
					"evictionHard": map[string]interface{}{
						"memory.available": "500Mi",
					},
				},
				"machineConfigPoolSelector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"aro.openshift.io/limits": "",
					},
				},
			},
			"status": map[string]interface{}{
				"conditions": nil,
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("systemreserved.kubeletConfig() = %v, want %v", got, want)
		t.Error(cmp.Diff(got, want))
	}
}
