package vacuum

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func createEc2ServiceForRegion(region string) (*ec2.EC2, error) {
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return ec2.New(mySession), nil
}

func Volumes(regions ...string) error {
	for _, region := range regions {
		svc, err := createEc2ServiceForRegion(region)
		if err != nil {
			return err
		}
		result, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("status"),
					Values: []*string{aws.String(ec2.VolumeStateAvailable)},
				},
			}})

		for _, volume := range result.Volumes {
			_, err := svc.DeleteVolume(&ec2.DeleteVolumeInput{VolumeId: volume.VolumeId})
			if err != nil {
				return err
			}
		}

	}

	return nil
}
