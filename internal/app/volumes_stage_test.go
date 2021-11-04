package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/kevholditch/vacuum/internal/app/vacuum"

	"github.com/avast/retry-go"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type volumesTestStage struct {
	t *testing.T
}

func newVolumesTest(t *testing.T) (*volumesTestStage, *volumesTestStage, *volumesTestStage) {
	s := &volumesTestStage{t: t}
	return s, s, s
}

func (s *volumesTestStage) createEc2ServiceForRegion(region string) *ec2.EC2 {
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		assert.Fail(s.t, err.Error())
	}
	return ec2.New(mySession)
}

func (s *volumesTestStage) an_available_volume_exists_in_region(region string) *volumesTestStage {
	svc := s.createEc2ServiceForRegion(region)
	v, err := svc.CreateVolume(&ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(region + "a"),
		VolumeType:       aws.String("gp2"),
		Size:             aws.Int64(10),
	})
	if err != nil {
		assert.Fail(s.t, "could not create volume, error: %v", err)
	}
	assert.True(s.t, v.VolumeId != nil)

	err = retry.Do(
		func() error {
			result, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("status"),
						Values: []*string{aws.String(ec2.VolumeStateAvailable)},
					},
				}})
			if err != nil {
				return err
			}
			if len(result.Volumes) == 0 {
				return fmt.Errorf("volume not created yet")
			}

			return nil
		}, retry.Delay(500*time.Millisecond))

	assert.NoError(s.t, err)

	return s
}

func (s *volumesTestStage) volumes_are_vacuumed_in(region string) *volumesTestStage {
	err := vacuum.Volumes(region)
	assert.NoError(s.t, err)
	return s
}

func (s *volumesTestStage) there_should_be_no_available_volumes_in_the_region(region string) *volumesTestStage {
	svc := s.createEc2ServiceForRegion(region)

	result, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String(ec2.VolumeStateAvailable)},
			},
		}})
	require.NoError(s.t, err)
	assert.Equal(s.t, 0, len(result.Volumes), fmt.Sprintf("should not be any available volumes in region: %s, found %d", region, len(result.Volumes)))

	return s
}

func (s *volumesTestStage) there_should_be_available_volumes_in_the_region(region string) *volumesTestStage {
	svc := s.createEc2ServiceForRegion(region)

	result, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String(ec2.VolumeStateAvailable)},
			},
		}})
	require.NoError(s.t, err)
	assert.True(s.t, len(result.Volumes) > 0, fmt.Sprintf("there should be available volumes in the region: %s", region))

	return s
}
