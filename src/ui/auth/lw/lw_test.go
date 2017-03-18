package lw

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/pkg/oidc"
	"github.com/vmware/harbor/src/pkg/tokenverifier"
)

const (
	testDomain string = "test.local"
)

func TestLogin(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	client := oidc.NewMockClient(controller)
	verifier := tokenverifier.NewMockTokenVerifier(controller)
	auth := &Auth{initialized: true, lwClient: client, tokenVerifier: verifier, lwDomainName: testDomain}

	token := oidc.NewMockTokenResponse(controller)

	client.EXPECT().GetTokenByPasswordGrant(gomock.Eq("testuser@test.local"), gomock.Eq("testpassword")).
		Times(1).Return(token, nil)

	auth.Login("testuser", "testpassword")
}

func TestEnsureUserProjects(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	client := oidc.NewMockClient(controller)
	verifier := tokenverifier.NewMockTokenVerifier(controller)
	auth := &Auth{initialized: true, lwClient: client, tokenVerifier: verifier, lwDomainName: testDomain}

	token := oidc.NewMockTokenResponse(controller)
	token.EXPECT().GetAccessToken().Times(1).Return("some_access_token_string")
	u := models.User{Username: "testuser"}

	jwtToken := &tokenverifier.JWTToken{Groups: []string{}}
	verifier.EXPECT().Verify("some_access_token_string").Times(1).Return(jwtToken, nil)

	auth.ensureUserProjects(u, token)
}

func TestGetGroups(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	client := oidc.NewMockClient(controller)
	verifier := tokenverifier.NewMockTokenVerifier(controller)
	auth := &Auth{initialized: true, lwClient: client, tokenVerifier: verifier, lwDomainName: testDomain}

	token := oidc.NewMockTokenResponse(controller)
	token.EXPECT().GetAccessToken().Times(1).Return("some_access_token_string")

	jwtToken := &tokenverifier.JWTToken{
		Groups: []string{
			"test.Local\\\\Harbor_ProjectTestABC_Admin",
			"test.Local\\\\irrelevant_groupA",
			"test.local\\\\Users",
		}}
	verifier.EXPECT().Verify("some_access_token_string").Times(1).Return(jwtToken, nil)

	groupInfos := auth.getUserProjectsFromToken(token)

	assert.NotNil(t, groupInfos)
	assert.Equal(t, len(groupInfos), 1)
	assert.Equal(t, groupInfos[0], GroupRole{groupName: "TestABC", role: "Admin"})
}

/* TODO: mock dao
func TestAuthenticate(t *testing.T) {
	client := oidc.NewMockClient(controller)
	verifier := tokenverifier.NewMockTokenVerifier(controller)
	auth := &Auth{initialized: true, LWClient: client, TokenVerifier: verifier, LWDomainName: testDomain}

	m := models.AuthModel{UserName: "testuser", Password: "testpassword"}
	auth.Authenticate(m)
}
*/
