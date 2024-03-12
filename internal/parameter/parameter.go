package parameter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
)

func GetParameter(region string, parameterName string, withDecrypt bool, assumeRole string) (string, error) {
	var sess *session.Session
	var err error

	if assumeRole == "" {
		sess, err = createSession(region)
	} else {
		sess, err = createAssumeRoleSession(assumeRole, "aws-inject", 300, region)
	}
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

func createSession(region string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func createAssumeRoleSession(roleArn string, sessionName string, duration int64, region string) (*session.Session, error) {
	sess := session.Must(session.NewSession())
	stsClient := sts.New(sess)

	assumeRoleOutput, err := stsClient.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int64(duration),
	})
	if err != nil {
		return nil, err
	}

	accessKey := aws.StringValue(assumeRoleOutput.Credentials.AccessKeyId)
	secretKey := aws.StringValue(assumeRoleOutput.Credentials.SecretAccessKey)
	sessionToken := aws.StringValue(assumeRoleOutput.Credentials.SessionToken)

	credentials := credentials.NewStaticCredentials(accessKey, secretKey, sessionToken)

	sessWithTempCreds := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials,
		Region:      aws.String(region),
	}))

	return sessWithTempCreds, nil
}
