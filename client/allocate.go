package client

import (
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// AllocationResult contains a net.IP / Interface pair
type AllocationResult struct {
	*net.IP
	Interface Interface
}

// AllocateClient offers IP allocation on interfaces
type AllocateClient interface {
	DeallocateIP(ipToRelease *net.IP) error
}

type allocateClient struct {
	aws    *awsclient
	subnet SubnetsClient
}

// DeallocateIP releases an IP back to AWS
func (c *allocateClient) DeallocateIP(ipToRelease *net.IP) error {
	client, err := c.aws.newEC2()
	if err != nil {
		return err
	}
	interfaces, err := c.aws.GetInterfaces()
	if err != nil {
		return err
	}

	for _, intf := range interfaces {
		for _, ip := range intf.IPv4s {
			if ipToRelease.Equal(ip) {
				request := ec2.UnassignPrivateIpAddressesInput{}
				request.SetNetworkInterfaceId(intf.ID)
				strIP := ipToRelease.String()
				request.SetPrivateIpAddresses([]*string{&strIP})
				_, err = client.UnassignPrivateIpAddresses(&request)
				//The primary IP address of an interface cannot be unassigned and continuing the loop
				if !strings.Contains(err.Error(), "InvalidParameterValue") {
					return nil
				}
			}
		}
	}
	return err
}
