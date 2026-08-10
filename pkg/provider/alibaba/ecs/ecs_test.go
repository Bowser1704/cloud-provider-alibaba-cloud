package ecs

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/stretchr/testify/assert"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/model"
	"k8s.io/cloud-provider-alibaba-cloud/pkg/provider/alibaba/base"
)

func NewECSClient() (*ecs.Client, error) {
	var ak, sk, regionId string
	if ak == "" || sk == "" {
		return nil, fmt.Errorf("ak or sk is empty")
	}
	return ecs.NewClientWithAccessKey(regionId, ak, sk)
}

func TestEcsProvider_ListInstances(t *testing.T) {
	client, err := NewECSClient()
	if err != nil {
		t.Skip("fail to create ecs client, skip")
		return
	}
	ids := []string{
		"cn-hangzhou.i-xxxx",
	}

	ecsProvider := NewECSProvider(&base.ClientMgr{
		Meta: nil,
		ECS:  client,
		VPC:  nil,
		SLB:  nil,
		PVTZ: nil,
	})

	cloudNodes, err := ecsProvider.ListInstances(context.TODO(), ids)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	for _, id := range ids {
		c, ok := cloudNodes[id]
		if !ok || c == nil {
			t.Errorf("cannot find %s ecs from cloud", id)
		}
		t.Logf("%s, instance address: %s", id, c.Addresses)
	}

	t.Logf("ListInstances test successfully")
}

func TestCollectNetworkInterfaceAddresses(t *testing.T) {
	networkInterface := ecs.NetworkInterfaceSet{
		NetworkInterfaceId: "eni-prefix",
		PrivateIpSets: ecs.PrivateIpSetsInDescribeNetworkInterfaces{
			PrivateIpSet: []ecs.PrivateIpSet{{PrivateIpAddress: "10.0.0.2"}},
		},
		Ipv4PrefixSets: ecs.Ipv4PrefixSetsInDescribeNetworkInterfaces{
			Ipv4PrefixSet: []ecs.Ipv4PrefixSet{{Ipv4Prefix: "10.0.0.16/28"}},
		},
	}

	result := map[string]string{}
	collectNetworkInterfaceAddresses(result, networkInterface, model.IPv4)

	assert.Equal(t, map[string]string{
		"10.0.0.2":     "eni-prefix",
		"10.0.0.16/28": "eni-prefix",
	}, result)
}
