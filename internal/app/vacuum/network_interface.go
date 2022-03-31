package vacuum

import (
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type networkInterfaceVacuumer struct{}

type networkInterfaceResource struct {
	id *string
}

func (n *networkInterfaceResource) ID() *string {
	return n.id
}

func NetworkInterfaces() Vacuumer {
	return &networkInterfaceVacuumer{}
}

func (v *networkInterfaceVacuumer) Type() string {
	return "network interfaces"
}

func (v *networkInterfaceVacuumer) Identify(region Region) (Resources, error) {
	svc, err := createEc2ServiceForRegion(region)
	response, err := svc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String(ec2.NetworkInterfaceStatusAvailable)},
			},
		}})
	if err != nil {
		return nil, err
	}
	result := &resources{
		resources: []Resource{},
		region:    region,
	}
	for _, n := range response.NetworkInterfaces {
		result.resources = append(result.resources, &networkInterfaceResource{id: n.NetworkInterfaceId})
	}

	return result, nil

}

func (v *networkInterfaceVacuumer) Clean(resources Resources, cleaned func(amount int)) error {
	svc, err := createEc2ServiceForRegion(resources.Region())
	if err != nil {
		return err
	}

	workers := 10
	workChan := make(chan *string)
	var amountCleaned int32

	for i := 0; i < workers; i++ {
		go func(workerNumber int) {
			for {
				id, ok := <-workChan
				if !ok {
					return
				}
				svc.DeleteNetworkInterface(&ec2.DeleteNetworkInterfaceInput{NetworkInterfaceId: id})
				cleaned(int(atomic.AddInt32(&amountCleaned, 1)))

			}

		}(i)
	}
	for _, resource := range resources.Resources() {
		workChan <- resource.ID()
	}
	close(workChan)

	return nil
}
