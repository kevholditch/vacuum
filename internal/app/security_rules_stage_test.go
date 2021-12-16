package app

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kevholditch/vacuum/internal/app/vacuum"
	"github.com/stretchr/testify/assert"
)

type securityRulesTestStage struct {
	t                *testing.T
	regions          []vacuum.Region
	resources        vacuum.Resources
	securityGroupdId *string
	vpcId            *string
}

func newSecurityRulesTest(t *testing.T) (*securityRulesTestStage, *securityRulesTestStage, *securityRulesTestStage, func()) {
	s := &securityRulesTestStage{t: t, regions: []vacuum.Region{}}
	return s, s, s, s.clean
}

func (s *securityRulesTestStage) createEc2ServiceForRegion(region string) *ec2.EC2 {
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		assert.Fail(s.t, err.Error())
	}
	return ec2.New(mySession)
}

func (s *securityRulesTestStage) createTestSecurityGroupForRegion(region string) *securityRulesTestStage {
	s.regions = append(s.regions, vacuum.Region(region))
	vacuumTestName := fmt.Sprintf("vacuum-test-%d", rand.Int())

	svc := s.createEc2ServiceForRegion(region)
	createVpcOutput, err := svc.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.1/16"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("vpc"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(vacuumTestName),
					},
				},
			},
		},
	})
	if err != nil {
		assert.Fail(s.t, err.Error())
	}

	s.vpcId = createVpcOutput.Vpc.VpcId
	createSecurityGroupOutput, err := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(vacuumTestName),
		Description: aws.String(vacuumTestName),
		VpcId:       aws.String(*s.vpcId),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidVpcID.NotFound":
				assert.Fail(s.t, fmt.Errorf("Unable to find VPC with ID %q.", *s.vpcId).Error())
			case "InvalidGroup.Duplicate":
				assert.Fail(s.t, fmt.Errorf("Security group %v already exists.", vacuumTestName).Error())
			}
		}
		assert.Fail(s.t, fmt.Errorf("Unable to create security group %q, %v", vacuumTestName, err).Error())
	}

	s.securityGroupdId = createSecurityGroupOutput.GroupId

	return s
}

func (s *securityRulesTestStage) createTestMatchingSecurityGroupRule(region string, egress bool) *securityRulesTestStage {
	testDescriptions := []string{
		"kubernetes.io/rule/nlb/client=a7a67774ec364454892df58fd6be9775",
		"kubernetes.io/rule/nlb/health=a6ba247c8d1e94b0691e80c29a769193",
		"kubernetes.io/rule/nlb/health=d1e94b069d474e4a247c8f4454dnlg54",
	}

	s.createTestSecurityGroupRule(region, testDescriptions[rand.Intn(len(testDescriptions))], egress)
	return s
}

func (s *securityRulesTestStage) createTestNonMatchingSecurityGroupRule(region string, egress bool) *securityRulesTestStage {
	testDescriptions := []string{
		"kubernetes.io/rule/nlb/mtu",
		"some-random=label",
		"cccccc uvtfdthketlucl kfjhhcdnlgdg dnrbirfvl lhn",
	}

	s.createTestSecurityGroupRule(region, testDescriptions[rand.Intn(len(testDescriptions))], egress)
	return s
}

func (s *securityRulesTestStage) createTestSecurityGroupRule(region string, description string, egress bool) *securityRulesTestStage {
	svc := s.createEc2ServiceForRegion(region)

	randomPort := 30000 + rand.Intn(5000)
	testPermissions := []*ec2.IpPermission{(&ec2.IpPermission{}).
		SetIpProtocol("tcp").
		SetFromPort(*aws.Int64(int64(randomPort))).
		SetToPort(*aws.Int64(int64(randomPort))).
		SetIpRanges([]*ec2.IpRange{
			{
				CidrIp:      aws.String("0.0.0.0/0"),
				Description: aws.String(description),
			},
		})}

	var err error
	if egress {
		_, err = svc.AuthorizeSecurityGroupEgress(&ec2.AuthorizeSecurityGroupEgressInput{
			IpPermissions: testPermissions,
			GroupId:       aws.String(*s.securityGroupdId),
		})
	} else {
		_, err = svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			IpPermissions: testPermissions,
			GroupId:       aws.String(*s.securityGroupdId),
		})
	}
	if err != nil {
		assert.Fail(s.t, err.Error())
	}

	return s
}

func (s *securityRulesTestStage) clean() {
	for _, region := range s.regions {
		svc := s.createEc2ServiceForRegion(string(region))
		_, err := svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{GroupId: s.securityGroupdId})
		if err != nil {
			assert.Fail(s.t, err.Error())
		}
		_, err = svc.DeleteVpc(&ec2.DeleteVpcInput{VpcId: s.vpcId})
		if err != nil {
			assert.Fail(s.t, err.Error())
		}
	}
}

func (s *securityRulesTestStage) a_security_group_with_no_rules_exists_in_region(region string) *securityRulesTestStage {
	return s.createTestSecurityGroupForRegion(region)
}

func (s *securityRulesTestStage) a_security_group_with_no_matching_rules_exists_in_region(region string) *securityRulesTestStage {
	return s.createTestSecurityGroupForRegion(region).
		createTestNonMatchingSecurityGroupRule(region, true).
		createTestNonMatchingSecurityGroupRule(region, false)
}

func (s *securityRulesTestStage) a_security_group_with_six_rules_and_three_matches_exists_in_region(region string) *securityRulesTestStage {
	return s.createTestSecurityGroupForRegion(region).
		createTestNonMatchingSecurityGroupRule(region, true).
		createTestNonMatchingSecurityGroupRule(region, false).
		createTestNonMatchingSecurityGroupRule(region, true).
		createTestMatchingSecurityGroupRule(region, false).
		createTestMatchingSecurityGroupRule(region, true).
		createTestMatchingSecurityGroupRule(region, false)
}

func (s *securityRulesTestStage) a_security_group_with_six_rules_and_six_matches_exists_in_region(region string) *securityRulesTestStage {
	return s.createTestSecurityGroupForRegion(region).
		createTestMatchingSecurityGroupRule(region, true).
		createTestMatchingSecurityGroupRule(region, false).
		createTestMatchingSecurityGroupRule(region, true).
		createTestMatchingSecurityGroupRule(region, false).
		createTestMatchingSecurityGroupRule(region, true).
		createTestMatchingSecurityGroupRule(region, false)
}

func (s *securityRulesTestStage) security_rules_are_identified_in(region string) *securityRulesTestStage {
	var err error
	s.resources, err = vacuum.SecurityRules().Identify(vacuum.Region(region))
	assert.NoError(s.t, err)

	return s
}

func (s *securityRulesTestStage) security_rules_are_vacuumed_in(region string) *securityRulesTestStage {
	resources, err := vacuum.SecurityRules().Identify(vacuum.Region(region))
	assert.NoError(s.t, err)

	err = vacuum.SecurityRules().Clean(resources, func(amount int) {})
	assert.NoError(s.t, err)

	return s
}

func (s *securityRulesTestStage) expected_number_equals_resources_length(expected int) {
	assert.Equal(s.t, expected, len(s.resources.Resources()))
}

func (s *securityRulesTestStage) there_should_be_no_security_rules_identified() *securityRulesTestStage {
	s.expected_number_equals_resources_length(0)
	return s
}

func (s *securityRulesTestStage) there_should_be_three_security_rules_identified() *securityRulesTestStage {
	s.expected_number_equals_resources_length(3)
	return s
}

func (s *securityRulesTestStage) there_should_be_six_security_rules_identified() *securityRulesTestStage {
	s.expected_number_equals_resources_length(6)
	return s
}
