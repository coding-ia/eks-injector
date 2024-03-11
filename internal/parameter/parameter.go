package parameter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetParameter(region string, parameterName string, withDecrypt bool) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", err
	}

	svc := ssm.New(sess)

	result, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(withDecrypt), // Set to true if the parameter is encrypted
	})
	if err != nil {
		return "", err
	}

	return *result.Parameter.Value, nil
}
