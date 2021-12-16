package app

import (
	"testing"
)

func TestIdentifySecurityRulesNoRules(t *testing.T) {
	given, when, then, clean := newSecurityRulesTest(t)
	defer clean()

	given.
		a_security_group_with_no_rules_exists_in_region("eu-west-1")

	when.
		security_rules_are_identified_in("eu-west-1")

	then.
		there_should_be_no_security_rules_identified()

}

func TestIdentifySecurityRulesNoMatches(t *testing.T) {
	given, when, then, clean := newSecurityRulesTest(t)
	defer clean()

	given.
		a_security_group_with_no_matching_rules_exists_in_region("eu-west-1")

	when.
		security_rules_are_identified_in("eu-west-1")

	then.
		there_should_be_no_security_rules_identified()

}

func TestIdentifySecurityRulesSixRulesAndThreeMatches(t *testing.T) {
	given, when, then, clean := newSecurityRulesTest(t)
	defer clean()

	given.
		a_security_group_with_six_rules_and_three_matches_exists_in_region("eu-west-1")

	when.
		security_rules_are_identified_in("eu-west-1")

	then.
		there_should_be_three_security_rules_identified()

}

func TestIdentifySecurityRulesSixRulesAndSixMatches(t *testing.T) {
	given, when, then, clean := newSecurityRulesTest(t)
	defer clean()

	given.
		a_security_group_with_six_rules_and_six_matches_exists_in_region("eu-west-1")

	when.
		security_rules_are_identified_in("eu-west-1")

	then.
		there_should_be_six_security_rules_identified()

}

func TestRemoveSecurityRulesSixRulesAndNoMatches(t *testing.T) {
	given, when, then, clean := newSecurityRulesTest(t)
	defer clean()

	given.
		a_security_group_with_no_matching_rules_exists_in_region("eu-west-1")

	when.
		security_rules_are_vacuumed_in("eu-west-1")

	then.
		security_rules_are_identified_in("eu-west-1").
		there_should_be_no_security_rules_identified()

}

func TestRemoveSecurityRulesSixRulesAndThreeMatches(t *testing.T) {
	given, when, then, clean := newSecurityRulesTest(t)
	defer clean()

	given.
		a_security_group_with_six_rules_and_three_matches_exists_in_region("eu-west-1")

	when.
		security_rules_are_vacuumed_in("eu-west-1")

	then.
		security_rules_are_identified_in("eu-west-1").
		there_should_be_no_security_rules_identified()

}

func TestRemoveSecurityRulesSixRulesAndSixMatches(t *testing.T) {
	given, when, then, clean := newSecurityRulesTest(t)
	defer clean()

	given.
		a_security_group_with_six_rules_and_six_matches_exists_in_region("eu-west-1")

	when.
		security_rules_are_vacuumed_in("eu-west-1")

	then.
		security_rules_are_identified_in("eu-west-1").
		there_should_be_no_security_rules_identified()

}
