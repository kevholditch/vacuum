package app

import (
	"fmt"
	"strings"
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
	t         *testing.T
	regions   map[vacuum.Region][]string
	resources vacuum.Resources
}

func newVolumesTest(t *testing.T) (*volumesTestStage, *volumesTestStage, *volumesTestStage, func()) {
	s := &volumesTestStage{t: t, regions: map[vacuum.Region][]string{}}
	return s, s, s, s.clean
}

func (s *volumesTestStage) clean() {
	for region := range s.regions {
		resources, _ := vacuum.Volumes().Identify(region)
		_ = vacuum.Volumes().Clean(resources, func(amount int) {})
	}
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

func (s *volumesTestStage) three_available_volumes_exist_in(region string) *volumesTestStage {
	return s.create_x_volumes_in_region(region, 3)
}

func (s *volumesTestStage) create_x_volumes_in_region(region string, amountOfVolumes int) *volumesTestStage {

	svc := s.createEc2ServiceForRegion(region)
	for i := 0; i < amountOfVolumes; i++ {

		v, err := svc.CreateVolume(&ec2.CreateVolumeInput{
			AvailabilityZone: aws.String(region + "a"),
			VolumeType:       aws.String("gp2"),
			Size:             aws.Int64(10),
		})
		if err != nil {
			assert.Fail(s.t, "could not create volume, error: %v", err)
		}
		assert.True(s.t, v.VolumeId != nil)

		volumes, ok := s.regions[vacuum.Region(region)]
		if !ok {
			s.regions[vacuum.Region(region)] = []string{*v.VolumeId}
		} else {
			s.regions[vacuum.Region(region)] = append(volumes, *v.VolumeId)
		}
	}

	err := retry.Do(
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
			if len(result.Volumes) != amountOfVolumes {
				return fmt.Errorf("volume(s) not created yet")
			}

			return nil
		}, retry.Delay(500*time.Millisecond))

	assert.NoError(s.t, err)

	return s
}

func (s *volumesTestStage) an_available_volume_exists_in_region(region string) *volumesTestStage {
	return s.create_x_volumes_in_region(region, 1)
}

func (s *volumesTestStage) the_available_volumes_are_identified(region string) *volumesTestStage {
	var err error
	s.resources, err = vacuum.Volumes().Identify(vacuum.Region(region))
	assert.NoError(s.t, err)

	return s
}

func (s *volumesTestStage) volumes_are_vacuumed_in(region string) *volumesTestStage {
	resources, err := vacuum.Volumes().Identify(vacuum.Region(region))
	assert.NoError(s.t, err)

	err = vacuum.Volumes().Clean(resources, func(amount int) {})
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

func (s *volumesTestStage) there_should_be_three_available_volumes_identified_in(region string) *volumesTestStage {

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
	volumeIds := s.regions[vacuum.Region(region)]

	for _, resource := range s.resources.Resources() {
		found := findElement(*resource.ID(), volumeIds)
		if !found {
			assert.Fail(s.t, fmt.Sprintf("failed to identify resource with id: %s", *resource.ID()))
		}
	}

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
