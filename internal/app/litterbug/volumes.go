package litterbug

import (
	"fmt"
	"time"

	"github.com/liamg/tml"

	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type bug struct {
	e               *ec2.EC2
	region          string
	volumesToCreate int
}

func InRegion(region string) (*bug, error) {
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return &bug{e: ec2.New(mySession), region: region}, nil
}

func (b *bug) CreateAvailableVolumes(amount int) *bug {
	b.volumesToCreate = amount
	return b
}

func (b *bug) Litter() error {
	tml.Printf("[%s]\n", b.region)
	for i := 0; i < b.volumesToCreate; i++ {

		v, err := b.e.CreateVolume(&ec2.CreateVolumeInput{
			AvailabilityZone: aws.String(b.region + "a"),
			VolumeType:       aws.String("gp2"),
			Size:             aws.Int64(10),
		})
		if err != nil {
			return err
		}

		tml.Printf("\tcreated volume: %s\n", *v.VolumeId)
	}

	err := retry.Do(
		func() error {
			result, err := b.e.DescribeVolumes(&ec2.DescribeVolumesInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("status"),
						Values: []*string{aws.String(ec2.VolumeStateAvailable)},
					},
				}})
			if err != nil {
				return err
			}
			if len(result.Volumes) != b.volumesToCreate {
				return fmt.Errorf("volume(s) not created yet")
			}

			return nil
		}, retry.Delay(500*time.Millisecond))

	return err
}
