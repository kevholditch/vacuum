package vacuum

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func createEc2ServiceForRegion(region Region) (*ec2.EC2, error) {
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String(string(region)),
	})
	if err != nil {
		return nil, err
	}
	return ec2.New(mySession), nil
}

type volumerVacuumer struct{}

type volumeResource struct {
	id *string
}

func (v *volumeResource) ID() *string {
	return v.id
}

func Volumes() Vacuumer {
	return &volumerVacuumer{}
}

func (v *volumerVacuumer) Type() string {
	return "volumes"
}

func (v *volumerVacuumer) Identify(region Region) (Resources, error) {
	svc, err := createEc2ServiceForRegion(region)
	if err != nil {
		return nil, err
	}

	response, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String(ec2.VolumeStateAvailable)},
			},
		}})

	if err != nil {
		return nil, err
	}

	result := &resources{
		resources: []Resource{},
		region:    region,
	}
	for _, v := range response.Volumes {
		result.resources = append(result.resources, &volumeResource{id: v.VolumeId})
	}

	return result, nil
}

func (v *volumerVacuumer) Clean(resources Resources, cleaned func(amount int)) error {
	svc, err := createEc2ServiceForRegion(resources.Region())
	if err != nil {
		return err
	}

	for i, resource := range resources.Resources() {
		_, err := svc.DeleteVolume(&ec2.DeleteVolumeInput{VolumeId: resource.ID()})
		if err != nil {
			return err
		}
		cleaned(i)
	}

	return nil
}
