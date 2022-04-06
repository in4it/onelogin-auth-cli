package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func AssumeRoleWithSAML(accountID string, role string, assertionPayload string) (*sts.AssumeRoleWithSAMLOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return nil, err
	}

	svc := sts.New(sess)

	roleToAssumeArn := "arn:aws:iam::" + accountID + ":role/" + role
	result, err := svc.AssumeRoleWithSAML(&sts.AssumeRoleWithSAMLInput{
		RoleArn:         &roleToAssumeArn,
		DurationSeconds: aws.Int64(3600),
		PrincipalArn:    aws.String("arn:aws:iam::" + accountID + ":saml-provider/" + role),
		SAMLAssertion:   aws.String(assertionPayload),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
