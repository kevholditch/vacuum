package app

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/avast/retry-go"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"

	"github.com/kevholditch/vacuum/internal/app/vacuum"
)

type networkInterfacesTestStage struct {
	t         *testing.T
	regions   map[vacuum.Region][]string
	resources vacuum.Resources
	subnetId  *string
	vpcId     *string
}

func newNetworkInterfacesTest(t *testing.T) (*networkInterfacesTestStage, *networkInterfacesTestStage, *networkInterfacesTestStage, func()) {
	s := &networkInterfacesTestStage{t: t, regions: map[vacuum.Region][]string{}}
	return s, s, s, s.clean
}

func (s *networkInterfacesTestStage) createTestSubnet(region string) *networkInterfacesTestStage {
	svc := s.createEc2ServiceForRegion(region)
	createVpcOutput, err := svc.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.1/16"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("vpc"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(fmt.Sprintf("vacuum-test-%d", rand.Int())),
					},
				},
			},
		},
	})
	if err != nil {
		s.t.Fail()
	}

	s.vpcId = createVpcOutput.Vpc.VpcId

	createSubnetOutput, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
		AvailabilityZone: aws.String(region + "a"),
		CidrBlock:        aws.String("10.0.0.1/24"),
		VpcId:            createVpcOutput.Vpc.VpcId,
	})
	if err != nil {
		s.t.Fail()
	}

	s.subnetId = createSubnetOutput.Subnet.SubnetId

	return s
}

func (s *networkInterfacesTestStage) clean() {
	for region := range s.regions {
		resources, _ := vacuum.NetworkInterfaces().Identify(region)
		_ = vacuum.NetworkInterfaces().Clean(resources)
		svc := s.createEc2ServiceForRegion(string(region))
		svc.DeleteSubnet(&ec2.DeleteSubnetInput{SubnetId: s.subnetId})
		svc.DeleteVpc(&ec2.DeleteVpcInput{VpcId: s.vpcId})
	}

}

func (s *networkInterfacesTestStage) createEc2ServiceForRegion(region string) *ec2.EC2 {
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		assert.Fail(s.t, err.Error())
	}
	return ec2.New(mySession)
}

func (s *networkInterfacesTestStage) an_available_network_interface_exists_in_region(region string) *networkInterfacesTestStage {
	return s.createTestSubnet(region).
		create_x_network_interfaces_in_region(region, 1)
}
func (s *networkInterfacesTestStage) three_available_network_interfaces_exist_in(region string) *networkInterfacesTestStage {
	return s.createTestSubnet(region).
		create_x_network_interfaces_in_region(region, 3)
}

func (s *networkInterfacesTestStage) network_interfaces_are_vacuumed_in(region string) *networkInterfacesTestStage {

	resources, err := vacuum.NetworkInterfaces().Identify(vacuum.Region(region))
	assert.NoError(s.t, err)

	err = vacuum.NetworkInterfaces().Clean(resources)
	assert.NoError(s.t, err)

	return s
}

func (s *networkInterfacesTestStage) create_x_network_interfaces_in_region(region string, amountOfNetworkInterfaces int) *networkInterfacesTestStage {

	svc := s.createEc2ServiceForRegion(region)

	for i := 0; i < amountOfNetworkInterfaces; i++ {

		n, err := svc.CreateNetworkInterface(&ec2.CreateNetworkInterfaceInput{SubnetId: s.subnetId})
		if err != nil {
			assert.Fail(s.t, "could not create network interface, error: %v", err)
		}
		assert.True(s.t, n.NetworkInterface.NetworkInterfaceId != nil)

		networkInterfaces, ok := s.regions[vacuum.Region(region)]
		if !ok {
			s.regions[vacuum.Region(region)] = []string{*n.NetworkInterface.NetworkInterfaceId}
		} else {
			s.regions[vacuum.Region(region)] = append(networkInterfaces, *n.NetworkInterface.NetworkInterfaceId)
		}
	}

	err := retry.Do(
		func() error {
			result, err := svc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("status"),
						Values: []*string{aws.String(ec2.NetworkInterfaceStatusAvailable)},
					},
				}})
			if err != nil {
				return err
			}
			if len(result.NetworkInterfaces) != amountOfNetworkInterfaces {
				return fmt.Errorf("network interface(s) not created yet")
			}

			return nil
		}, retry.Delay(500*time.Millisecond))

	assert.NoError(s.t, err)

	return s
}

func (s *networkInterfacesTestStage) the_available_network_interfaces_are_identified(region string) *networkInterfacesTestStage {
	var err error
	s.resources, err = vacuum.NetworkInterfaces().Identify(vacuum.Region(region))
	assert.NoError(s.t, err)

	return s
}

func (s *networkInterfacesTestStage) there_should_be_three_available_network_interfaces_identified_in(region string) *networkInterfacesTestStage {

	assert.Equal(s.t, region, string(s.resources.Region()))
	assert.Equal(s.t, 3, len(s.resources.Resources()))

	findElement := func(search string, elements []string) bool {
		for _, elem := range elements {
			if strings.EqualFold(elem, search) {
				return true
			}
		}
		return false
	}
	ids := s.regions[vacuum.Region(region)]

	for _, resource := range s.resources.Resources() {
		found := findElement(*resource.ID(), ids)
		if !found {
			assert.Fail(s.t, fmt.Sprintf("failed to identify resource with id: %s", *resource.ID()))
		}
	}

	return s
}

func (s *networkInterfacesTestStage) there_should_be_no_available_network_interfaces_in_the_region(region string) *networkInterfacesTestStage {

	svc := s.createEc2ServiceForRegion(region)

	result, err := svc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String(ec2.NetworkInterfaceStatusAvailable)},
			},
		}})
	require.NoError(s.t, err)
	assert.Equal(s.t, 0, len(result.NetworkInterfaces), fmt.Sprintf("should not be any available network interfaces in region: %s, found %d", region, len(result.NetworkInterfaces)))

	return s
}
