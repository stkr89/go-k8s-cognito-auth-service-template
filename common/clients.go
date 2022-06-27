package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/costandusagereportservice"
)

func NewAWSCognitoClient() *cognito.CognitoIdentityProvider {
	conf := &aws.Config{
		Region: aws.String(costandusagereportservice.AWSRegionUsEast1),
	}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	return cognito.New(sess)
}
