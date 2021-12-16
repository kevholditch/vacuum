package vacuum

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func createELBServiceForRegion(region Region) (*elb.ELB, error) {
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String(string(region)),
	})
	if err != nil {
		return nil, err
	}
	return elb.New(mySession), nil
}

type securityRuleVacuumer struct{}

type securityRuleResource struct {
	id *string
}

func (s *securityRuleResource) ID() *string {
	return s.id
}

func SecurityRules() Vacuumer {
	return &securityRuleVacuumer{}
}

func (s *securityRuleVacuumer) Type() string {
	return "security rules"
}

func (s *securityRuleVacuumer) Identify(region Region) (Resources, error) {
	svc, err := createEc2ServiceForRegion(region)
	if err != nil {
		return nil, err
	}

	response, err := svc.DescribeSecurityGroupRules(&ec2.DescribeSecurityGroupRulesInput{})
	if err != nil {
		return nil, err
	}

	result := &resources{
		resources: []Resource{},
		region:    region,
	}
	for _, rule := range response.SecurityGroupRules {
		removable, err := securityGroupRuleCanBeRemoved(rule, region)
		if err != nil {
			return nil, err
		}
		if !removable {
			continue
		}
		result.resources = append(result.resources, &securityRuleResource{id: rule.SecurityGroupRuleId})
	}

	return result, nil
}

func securityGroupRuleCanBeRemoved(rule *ec2.SecurityGroupRule, region Region) (bool, error) {
	k8sManaged, nlbID := securityGroupRuleCreatedByK8sManagedNLB(rule)
	if k8sManaged {
		nlbExists, err := nlbExists(nlbID, region)
		if err != nil {
			return false, err
		}
		if nlbExists {
			return false, nil
		}

		return true, nil
	}

	return false, nil
}

func securityGroupRuleCreatedByK8sManagedNLB(rule *ec2.SecurityGroupRule) (bool, *string) {
	// Example descriptions we are expecting to be able to handle:
	// kubernetes.io/rule/nlb/client=a7a67774ec364454892df58fd6be9775
	// kubernetes.io/rule/nlb/health=a6ba247c8d1e94b0691e80c29a769193
	if rule.Description == nil || !strings.HasPrefix(*rule.Description, "kubernetes.io/rule/nlb") {
		return false, nil
	}

	splitK8sDescription := strings.Split(*rule.Description, "=")
	if len(splitK8sDescription) == 1 {
		return false, nil
	}

	nlb := splitK8sDescription[len(splitK8sDescription)-1]
	return true, &nlb
}

func nlbExists(nlbID *string, region Region) (bool, error) {
	svc, err := createELBServiceForRegion(region)
	if err != nil {
		return false, nil
	}
	_, err = svc.DescribeLoadBalancerAttributes(&elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: nlbID,
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == elbv2.ErrCodeLoadBalancerNotFoundException {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}

func (s *securityRuleVacuumer) Clean(resources Resources) error {
	svc, err := createEc2ServiceForRegion(resources.Region())
	if err != nil {
		return err
	}

	var ruleIds []*string
	for _, resource := range resources.Resources() {
		ruleIds = append(ruleIds, resource.ID())
	}

	rulesResponse, err := svc.DescribeSecurityGroupRules(&ec2.DescribeSecurityGroupRulesInput{
		SecurityGroupRuleIds: ruleIds,
	})
	if err != nil {
		return err
	}

	egressRules, ingressRules := rulesSplitByEgressIngress(rulesResponse.SecurityGroupRules)

	for groupID, ruleIDs := range egressRules {
		_, err = svc.RevokeSecurityGroupEgress(&ec2.RevokeSecurityGroupEgressInput{
			GroupId:              &groupID,
			SecurityGroupRuleIds: ruleIDs,
		})
		if err != nil {
			return err
		}
	}

	for groupID, ruleIDs := range ingressRules {
		_, err = svc.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
			GroupId:              &groupID,
			SecurityGroupRuleIds: ruleIDs,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func rulesSplitByEgressIngress(rules []*ec2.SecurityGroupRule) (egressRules map[string][]*string, ingressRules map[string][]*string) {
	egressRules = make(map[string][]*string)
	ingressRules = make(map[string][]*string)

	for _, rule := range rules {
		if *rule.IsEgress {
			egressRules[*rule.GroupId] = append(egressRules[*rule.GroupId], rule.SecurityGroupRuleId)
		} else {
			ingressRules[*rule.GroupId] = append(ingressRules[*rule.GroupId], rule.SecurityGroupRuleId)
		}
	}

	return egressRules, ingressRules
}
