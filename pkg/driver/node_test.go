/*
Copyright 2021 Juicedata Inc

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

package driver

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/juicedata/juicefs-csi-driver/pkg/juicefs"
	"github.com/juicedata/juicefs-csi-driver/pkg/juicefs/mocks"
	podmount "github.com/juicedata/juicefs-csi-driver/pkg/juicefs/mount"
	"k8s.io/utils/mount"
	"os"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/container-storage-interface/spec/lib/go/csi"
	k8s "github.com/juicedata/juicefs-csi-driver/pkg/juicefs/k8sclient"
	. "github.com/smartystreets/goconvey/convey"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNodePublishVolume(t *testing.T) {
	stdVolCap := &csi.VolumeCapability{
		AccessType: &csi.VolumeCapability_Mount{
			Mount: &csi.VolumeCapability_MountVolume{},
		},
		AccessMode: &csi.VolumeCapability_AccessMode{
			Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		},
	}
	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "success normal",
			testFunc: func(t *testing.T) {
				Convey("Test NodePublishVolume", t, func() {
					Convey("test normal", func() {
						volumeId := "vol-test"
						subPath := "/subPath"
						targetPath := "/test/path"
						bindSource := "/test/path"
						volumeCtx := map[string]string{"subPath": subPath}
						secret := map[string]string{"a": "b"}

						patch1 := ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
							return nil
						})
						defer patch1.Reset()

						mockCtl := gomock.NewController(t)
						defer mockCtl.Finish()

						mockJfs := mocks.NewMockJfs(mockCtl)
						mockJfs.EXPECT().CreateVol(volumeId, subPath).Return(bindSource, nil)
						mockJuicefs := mocks.NewMockInterface(mockCtl)
						mockJuicefs.EXPECT().Mount(bindSource, targetPath, fsTypeNone, []string{"bind"}).Return(nil)
						mockJuicefs.EXPECT().JfsMount(volumeId, targetPath, secret, volumeCtx, []string{"ro"}, true).Return(mockJfs, nil)

						juicefsDriver := &nodeService{
							juicefs:   mockJuicefs,
							nodeID:    "fake_node_id",
							k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
						}

						req := &csi.NodePublishVolumeRequest{
							VolumeId:         volumeId,
							TargetPath:       targetPath,
							VolumeCapability: stdVolCap,
							Readonly:         true,
							Secrets:          secret,
							VolumeContext:    volumeCtx,
						}

						_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
						if err != nil {
							t.Fatalf("Expect no error but got: %v", err)
						}
					})
					Convey("test mountOptions in volumeAttributes", func() {
						volumeId := "vol-test"
						subPath := "/subPath"
						targetPath := "/test/path"
						bindSource := "/test/path"
						mountOptions := []string{"cache-dir=/cache"}
						volumeCtx := map[string]string{
							"subPath":      subPath,
							"mountOptions": "cache-dir=/cache",
						}
						secret := map[string]string{"a": "b"}

						patch1 := ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
							return nil
						})
						defer patch1.Reset()

						mockCtl := gomock.NewController(t)
						defer mockCtl.Finish()

						mockJfs := mocks.NewMockJfs(mockCtl)
						mockJfs.EXPECT().CreateVol(volumeId, subPath).Return(bindSource, nil)
						mockJuicefs := mocks.NewMockInterface(mockCtl)
						mockJuicefs.EXPECT().Mount(bindSource, targetPath, fsTypeNone, []string{"bind"}).Return(nil)
						mockJuicefs.EXPECT().JfsMount(volumeId, targetPath, secret, volumeCtx, mountOptions, true).Return(mockJfs, nil)

						juicefsDriver := &nodeService{
							juicefs:   mockJuicefs,
							nodeID:    "fake_node_id",
							k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
						}

						req := &csi.NodePublishVolumeRequest{
							VolumeId:         volumeId,
							TargetPath:       targetPath,
							VolumeCapability: stdVolCap,
							Readonly:         false,
							Secrets:          secret,
							VolumeContext:    volumeCtx,
						}

						_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
						if err != nil {
							t.Fatalf("Expect no error but got: %v", err)
						}
					})
					Convey("test mountOptions in spec", func() {
						volumeId := "vol-test"
						subPath := "/subPath"
						targetPath := "/test/path"
						bindSource := "/test/path"
						mountOptions := []string{"cache-dir=/cache"}
						volumeCtx := map[string]string{
							"subPath": subPath,
						}
						secret := map[string]string{"a": "b"}

						patch1 := ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
							return nil
						})
						defer patch1.Reset()

						mockCtl := gomock.NewController(t)
						defer mockCtl.Finish()

						mockJfs := mocks.NewMockJfs(mockCtl)
						mockJfs.EXPECT().CreateVol(volumeId, subPath).Return(bindSource, nil)
						mockJuicefs := mocks.NewMockInterface(mockCtl)
						mockJuicefs.EXPECT().Mount(bindSource, targetPath, fsTypeNone, []string{"bind"}).Return(nil)
						mockJuicefs.EXPECT().JfsMount(volumeId, targetPath, secret, volumeCtx, mountOptions, true).Return(mockJfs, nil)

						juicefsDriver := &nodeService{
							juicefs:   mockJuicefs,
							nodeID:    "fake_node_id",
							k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
						}

						stdVolCapWithMount := &csi.VolumeCapability{
							AccessType: &csi.VolumeCapability_Mount{
								Mount: &csi.VolumeCapability_MountVolume{
									MountFlags: mountOptions,
								},
							},
							AccessMode: &csi.VolumeCapability_AccessMode{
								Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
							},
						}
						req := &csi.NodePublishVolumeRequest{
							VolumeId:         volumeId,
							TargetPath:       targetPath,
							VolumeCapability: stdVolCapWithMount,
							Readonly:         false,
							Secrets:          secret,
							VolumeContext:    volumeCtx,
						}

						_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
						if err != nil {
							t.Fatalf("Expect no error but got: %v", err)
						}
					})
					Convey("test JfsMount err", func() {
						volumeId := "vol-test"
						subPath := "/subPath"
						targetPath := "/test/path"
						volumeCtx := map[string]string{"subPath": subPath}
						secret := map[string]string{"a": "b"}

						patch1 := ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
							return nil
						})
						defer patch1.Reset()

						mockCtl := gomock.NewController(t)
						defer mockCtl.Finish()

						mockJfs := mocks.NewMockJfs(mockCtl)
						mockJuicefs := mocks.NewMockInterface(mockCtl)
						mockJuicefs.EXPECT().JfsMount(volumeId, targetPath, secret, volumeCtx, []string{"ro"}, true).Return(mockJfs, errors.New("test"))

						juicefsDriver := &nodeService{
							juicefs:   mockJuicefs,
							nodeID:    "fake_node_id",
							k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
						}

						req := &csi.NodePublishVolumeRequest{
							VolumeId:         volumeId,
							TargetPath:       targetPath,
							VolumeCapability: stdVolCap,
							Readonly:         true,
							Secrets:          secret,
							VolumeContext:    volumeCtx,
						}

						_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
						if err == nil {
							t.Fatal("Expect error but got nil")
						}
					})
					Convey("test CreateVol err", func() {
						volumeId := "vol-test"
						subPath := "/subPath"
						targetPath := "/test/path"
						bindSource := "/test/path"
						volumeCtx := map[string]string{"subPath": subPath}
						secret := map[string]string{"a": "b"}

						patch1 := ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
							return nil
						})
						defer patch1.Reset()

						mockCtl := gomock.NewController(t)
						defer mockCtl.Finish()

						mockJfs := mocks.NewMockJfs(mockCtl)
						mockJfs.EXPECT().CreateVol(volumeId, subPath).Return(bindSource, errors.New("test"))
						mockJuicefs := mocks.NewMockInterface(mockCtl)
						mockJuicefs.EXPECT().JfsMount(volumeId, targetPath, secret, volumeCtx, []string{"ro"}, true).Return(mockJfs, nil)

						juicefsDriver := &nodeService{
							juicefs:   mockJuicefs,
							nodeID:    "fake_node_id",
							k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
						}

						req := &csi.NodePublishVolumeRequest{
							VolumeId:         volumeId,
							TargetPath:       targetPath,
							VolumeCapability: stdVolCap,
							Readonly:         true,
							Secrets:          secret,
							VolumeContext:    volumeCtx,
						}

						_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
						if err == nil {
							t.Fatal("Expect error but got nil")
						}
					})
					Convey("test Mount err", func() {
						volumeId := "vol-test"
						subPath := "/subPath"
						targetPath := "/test/path"
						bindSource := "/test/path"
						volumeCtx := map[string]string{"subPath": subPath}
						secret := map[string]string{"a": "b"}

						patch1 := ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
							return nil
						})
						defer patch1.Reset()

						mockCtl := gomock.NewController(t)
						defer mockCtl.Finish()

						mockJfs := mocks.NewMockJfs(mockCtl)
						mockJfs.EXPECT().CreateVol(volumeId, subPath).Return(bindSource, nil)
						mockJuicefs := mocks.NewMockInterface(mockCtl)
						mockJuicefs.EXPECT().Mount(bindSource, targetPath, fsTypeNone, []string{"bind"}).Return(errors.New("test"))
						mockJuicefs.EXPECT().JfsMount(volumeId, targetPath, secret, volumeCtx, []string{"ro"}, true).Return(mockJfs, nil)

						juicefsDriver := &nodeService{
							juicefs:   mockJuicefs,
							nodeID:    "fake_node_id",
							k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
						}

						req := &csi.NodePublishVolumeRequest{
							VolumeId:         volumeId,
							TargetPath:       targetPath,
							VolumeCapability: stdVolCap,
							Readonly:         true,
							Secrets:          secret,
							VolumeContext:    volumeCtx,
						}

						_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
						if err == nil {
							t.Fatal("Expect error but got nil")
						}
					})
					Convey("test MkdirAll err", func() {
						volumeId := "vol-test"
						subPath := "/subPath"
						targetPath := "/test/path"
						volumeCtx := map[string]string{"subPath": subPath}
						secret := map[string]string{"a": "b"}

						patch1 := ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error {
							return errors.New("test")
						})
						defer patch1.Reset()

						juicefsDriver := &nodeService{
							juicefs:   nil,
							nodeID:    "fake_node_id",
							k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
						}

						req := &csi.NodePublishVolumeRequest{
							VolumeId:         volumeId,
							TargetPath:       targetPath,
							VolumeCapability: stdVolCap,
							Readonly:         true,
							Secrets:          secret,
							VolumeContext:    volumeCtx,
						}

						_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
						if err == nil {
							t.Fatal("Expect error but got nil")
						}
					})
				})
			},
		},
		{
			name: "no target",
			testFunc: func(t *testing.T) {
				targetPath := ""

				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()
				mockJuicefs := mocks.NewMockInterface(mockCtl)
				juicefsDriver := &nodeService{
					juicefs:   mockJuicefs,
					nodeID:    "fake_node_id",
					k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
				}

				req := &csi.NodePublishVolumeRequest{
					TargetPath:       targetPath,
					VolumeCapability: stdVolCap,
				}

				_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
				if err == nil {
					t.Fatalf("Expect error but got nil")
				}
			},
		},
		{
			name: "no capability",
			testFunc: func(t *testing.T) {
				targetPath := "/test"

				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()
				mockJuicefs := mocks.NewMockInterface(mockCtl)
				juicefsDriver := &nodeService{
					juicefs:   mockJuicefs,
					nodeID:    "fake_node_id",
					k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
				}

				req := &csi.NodePublishVolumeRequest{
					TargetPath:       targetPath,
					VolumeCapability: nil,
				}

				_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
				if err == nil {
					t.Fatalf("Expect error but got nil")
				}
			},
		},
		{
			name: "invalid capability",
			testFunc: func(t *testing.T) {
				targetPath := "/test"

				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()
				mockJuicefs := mocks.NewMockInterface(mockCtl)
				juicefsDriver := &nodeService{
					juicefs:   mockJuicefs,
					nodeID:    "fake_node_id",
					k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
				}

				invalidVolumeCaps := &csi.VolumeCapability{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
					},
				}
				req := &csi.NodePublishVolumeRequest{
					TargetPath:       targetPath,
					VolumeCapability: invalidVolumeCaps,
				}

				_, err := juicefsDriver.NodePublishVolume(context.TODO(), req)
				if err == nil {
					t.Fatalf("Expect error but got nil")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestNodeUnpublishVolume(t *testing.T) {
	Convey("Test NodePublishVolume", t, func() {
		Convey("test normal", func() {
			targetPath := "/test/path"
			podMount := &podmount.PodMount{}
			patch2 := ApplyMethod(reflect.TypeOf(podMount), "JUmount", func(_ *podmount.PodMount, volumeId, target string) error {
				return nil
			})
			defer patch2.Reset()

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockJuicefs := mocks.NewMockInterface(mockCtl)
			mockJuicefs.EXPECT().JfsUnmount(targetPath).Return(nil)
			mockJuicefs.EXPECT().JfsCleanupMountPoint(targetPath).Return(nil)

			juicefsDriver := &nodeService{
				juicefs:   mockJuicefs,
				nodeID:    "fake_node_id",
				k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
			}

			req := &csi.NodeUnpublishVolumeRequest{
				TargetPath: targetPath,
				VolumeId:   "vol-test",
			}

			_, err := juicefsDriver.NodeUnpublishVolume(context.TODO(), req)
			if err != nil {
				t.Fatalf("Expect no error but got: %v", err)
			}
		})
		Convey("JfsUnmount err", func() {
			targetPath := "/test/path"
			podMount := &podmount.PodMount{}
			patch := ApplyMethod(reflect.TypeOf(podMount), "JUmount", func(_ *podmount.PodMount, volumeId, target string) error {
				return nil
			})
			defer patch.Reset()

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockJuicefs := mocks.NewMockInterface(mockCtl)
			mockJuicefs.EXPECT().JfsUnmount(targetPath).Return(errors.New("test"))

			juicefsDriver := &nodeService{
				juicefs:   mockJuicefs,
				nodeID:    "fake_node_id",
				k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
			}

			req := &csi.NodeUnpublishVolumeRequest{
				TargetPath: targetPath,
				VolumeId:   "vol-test",
			}

			_, err := juicefsDriver.NodeUnpublishVolume(context.TODO(), req)
			if err == nil {
				t.Fatal("Expect error but got nil")
			}
		})
		Convey("CleanupMountPoint err", func() {
			targetPath := "/test/path"
			podMount := &podmount.PodMount{}
			patch2 := ApplyMethod(reflect.TypeOf(podMount), "JUmount", func(_ *podmount.PodMount, volumeId, target string) error {
				return nil
			})
			defer patch2.Reset()

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockJuicefs := mocks.NewMockInterface(mockCtl)
			mockJuicefs.EXPECT().JfsUnmount(targetPath).Return(nil)
			mockJuicefs.EXPECT().JfsCleanupMountPoint(targetPath).Return(errors.New("test"))

			juicefsDriver := &nodeService{
				juicefs:   mockJuicefs,
				nodeID:    "fake_node_id",
				k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
			}

			req := &csi.NodeUnpublishVolumeRequest{
				TargetPath: targetPath,
				VolumeId:   "vol-test",
			}

			_, err := juicefsDriver.NodeUnpublishVolume(context.TODO(), req)
			if err == nil {
				t.Fatal("Expect error but got nil")
			}
		})
		Convey("JUmount err", func() {
			targetPath := "/test/path"
			podMount := &podmount.PodMount{}
			patch2 := ApplyMethod(reflect.TypeOf(podMount), "JUmount", func(_ *podmount.PodMount, volumeId, target string) error {
				return errors.New("test")
			})
			defer patch2.Reset()

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockJuicefs := mocks.NewMockInterface(mockCtl)
			mockJuicefs.EXPECT().JfsUnmount(targetPath).Return(nil)
			mockJuicefs.EXPECT().JfsCleanupMountPoint(targetPath).Return(nil)

			juicefsDriver := &nodeService{
				juicefs:   mockJuicefs,
				nodeID:    "fake_node_id",
				k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
			}

			req := &csi.NodeUnpublishVolumeRequest{
				TargetPath: targetPath,
				VolumeId:   "vol-test",
			}

			_, err := juicefsDriver.NodeUnpublishVolume(context.TODO(), req)
			if err == nil {
				t.Fatal("Expect error but got nil")
			}
		})
		Convey("nil target", func() {
			juicefsDriver := &nodeService{
				juicefs:   nil,
				nodeID:    "fake_node_id",
				k8sClient: &k8s.K8sClient{Interface: fake.NewSimpleClientset()},
			}

			req := &csi.NodeUnpublishVolumeRequest{
				TargetPath: "",
				VolumeId:   "vol-test",
			}

			_, err := juicefsDriver.NodeUnpublishVolume(context.TODO(), req)
			if err == nil {
				t.Fatal("Expect error but got nil")
			}
		})
	})
}

func Test_nodeService_NodeGetCapabilities(t *testing.T) {
	type fields struct {
		juicefs   juicefs.Interface
		nodeID    string
		k8sClient *k8s.K8sClient
	}
	type args struct {
		ctx context.Context
		req *csi.NodeGetCapabilitiesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *csi.NodeGetCapabilitiesResponse
		wantErr bool
	}{
		{
			name:    "test",
			fields:  fields{},
			args:    args{},
			want:    &csi.NodeGetCapabilitiesResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &nodeService{
				juicefs:   tt.fields.juicefs,
				nodeID:    tt.fields.nodeID,
				k8sClient: tt.fields.k8sClient,
			}
			got, err := d.NodeGetCapabilities(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeGetCapabilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeGetCapabilities() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeService_NodeGetInfo(t *testing.T) {
	type fields struct {
		juicefs   juicefs.Interface
		nodeID    string
		k8sClient *k8s.K8sClient
	}
	type args struct {
		ctx context.Context
		req *csi.NodeGetInfoRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *csi.NodeGetInfoResponse
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				nodeID: "test",
			},
			args:    args{},
			want:    &csi.NodeGetInfoResponse{NodeId: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &nodeService{
				juicefs:   tt.fields.juicefs,
				nodeID:    tt.fields.nodeID,
				k8sClient: tt.fields.k8sClient,
			}
			got, err := d.NodeGetInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeGetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeGetInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newNodeService(t *testing.T) {
	Convey("Test newNodeService", t, func() {
		Convey("normal", func() {
			patch1 := ApplyFunc(k8s.NewClient, func() (*k8s.K8sClient, error) {
				return nil, nil
			})
			defer patch1.Reset()
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockJuicefs := mocks.NewMockInterface(mockCtl)
			mockJuicefs.EXPECT().Version().Return([]byte(""), nil)

			patch2 := ApplyFunc(juicefs.NewJfsProvider, func(mounter *mount.SafeFormatAndMount) (juicefs.Interface, error) {
				return mockJuicefs, nil
			})
			defer patch2.Reset()
			_, err := newNodeService("test")
			So(err, ShouldBeNil)
		})
		Convey("NewClient err", func() {
			patch1 := ApplyFunc(k8s.NewClient, func() (*k8s.K8sClient, error) {
				return nil, errors.New("test")
			})
			defer patch1.Reset()
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockJuicefs := mocks.NewMockInterface(mockCtl)
			mockJuicefs.EXPECT().Version().Return([]byte(""), nil)

			patch2 := ApplyFunc(juicefs.NewJfsProvider, func(mounter *mount.SafeFormatAndMount) (juicefs.Interface, error) {
				return mockJuicefs, nil
			})
			defer patch2.Reset()
			_, err := newNodeService("test")
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_nodeService_NodeExpandVolume(t *testing.T) {
	type fields struct {
		juicefs   juicefs.Interface
		nodeID    string
		k8sClient *k8s.K8sClient
	}
	type args struct {
		ctx context.Context
		req *csi.NodeExpandVolumeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *csi.NodeExpandVolumeResponse
		wantErr bool
	}{
		{
			name:    "test",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &nodeService{
				juicefs:   tt.fields.juicefs,
				nodeID:    tt.fields.nodeID,
				k8sClient: tt.fields.k8sClient,
			}
			got, err := d.NodeExpandVolume(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeExpandVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeExpandVolume() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeService_NodeGetVolumeStats(t *testing.T) {
	type fields struct {
		juicefs   juicefs.Interface
		nodeID    string
		k8sClient *k8s.K8sClient
	}
	type args struct {
		ctx context.Context
		req *csi.NodeGetVolumeStatsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *csi.NodeGetVolumeStatsResponse
		wantErr bool
	}{
		{
			name:    "test",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &nodeService{
				juicefs:   tt.fields.juicefs,
				nodeID:    tt.fields.nodeID,
				k8sClient: tt.fields.k8sClient,
			}
			got, err := d.NodeGetVolumeStats(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeGetVolumeStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeGetVolumeStats() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeService_NodeStageVolume(t *testing.T) {
	type fields struct {
		juicefs   juicefs.Interface
		nodeID    string
		k8sClient *k8s.K8sClient
	}
	type args struct {
		ctx context.Context
		req *csi.NodeStageVolumeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *csi.NodeStageVolumeResponse
		wantErr bool
	}{
		{
			name:    "test",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &nodeService{
				juicefs:   tt.fields.juicefs,
				nodeID:    tt.fields.nodeID,
				k8sClient: tt.fields.k8sClient,
			}
			got, err := d.NodeStageVolume(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeStageVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeStageVolume() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeService_NodeUnstageVolume(t *testing.T) {
	type fields struct {
		juicefs   juicefs.Interface
		nodeID    string
		k8sClient *k8s.K8sClient
	}
	type args struct {
		ctx context.Context
		req *csi.NodeUnstageVolumeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *csi.NodeUnstageVolumeResponse
		wantErr bool
	}{
		{
			name:    "test",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &nodeService{
				juicefs:   tt.fields.juicefs,
				nodeID:    tt.fields.nodeID,
				k8sClient: tt.fields.k8sClient,
			}
			got, err := d.NodeUnstageVolume(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeUnstageVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeUnstageVolume() got = %v, want %v", got, tt.want)
			}
		})
	}
}
