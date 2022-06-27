package service

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/go-kit/kit/log"
	"github.com/stkr89/go-auth-service-template/common"
	"github.com/stkr89/go-auth-service-template/types"
	"os"
)

// AuthService interface
type AuthService interface {
	SignUp(ctx context.Context, request *types.SignUpRequest) (*types.SignUpResponse, error)
	SignIn(ctx context.Context, user *types.SignInRequest) (*types.SignInResponse, error)
}

type AuthServiceImpl struct {
	logger log.Logger
	client *cognito.CognitoIdentityProvider
}

func NewAuthServiceImpl() *AuthServiceImpl {
	return &AuthServiceImpl{
		logger: common.NewLogger(),
		client: common.NewAWSCognitoClient(),
	}
}

func (s AuthServiceImpl) SignIn(ctx context.Context, user *types.SignInRequest) (*types.SignInResponse, error) {
	authInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String(cognito.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(user.Email),
			"PASSWORD": aws.String(user.Password),
		},
		ClientId: s.getClientId(),
	}
	auth, err := s.client.InitiateAuth(authInput)
	if err != nil {
		s.logger.Log("message", "unable to signin user", "error", err)
		return nil, common.NewError(common.Unauthorized, "signin failed")
	}

	return &types.SignInResponse{
		AccessToken: *auth.AuthenticationResult.IdToken,
	}, nil
}

func (s AuthServiceImpl) SignUp(ctx context.Context, request *types.SignUpRequest) (*types.SignUpResponse, error) {
	awsUser := &cognito.SignUpInput{
		Username: aws.String(request.Email),
		Password: aws.String(request.Password),
		ClientId: s.getClientId(),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("custom:firstName"),
				Value: &request.FirstName,
			},
			{
				Name:  aws.String("custom:lastName"),
				Value: &request.LastName,
			},
		},
	}

	signUpOutput, err := s.client.SignUp(awsUser)
	if err != nil {
		s.logger.Log("message", "unable to signup user", "error", err)
		return nil, common.NewError(common.SomethingWentWrong, "signup failed")
	}

	s.logger.Log("message", "sign up successful", "email", request.Email)

	return &types.SignUpResponse{
		ID:        *signUpOutput.UserSub,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
	}, nil
}

func (s AuthServiceImpl) getClientId() *string {
	return aws.String(os.Getenv("AWS_COGNITO_CLIENT_ID"))
}
