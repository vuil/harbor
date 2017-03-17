package lw

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/pkg/oidc"
	"github.com/vmware/harbor/src/pkg/tokenverifier"
	"github.com/vmware/harbor/src/ui/config"
)

func getLWSetting() config.LightwaveSetting {
	return config.LightwaveSetting{
		DomainName:         "test.local",
		Endpoint:           "http://testserver",
		AdminUser:          "testAdmin",
		AdminPassword:      "testPassword!123",
		IgnoreCertificates: true,
	}
}

func TestLogin(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	token := oidc.NewMockTokenResponse(controller)

	client := oidc.NewMockClient(controller)
	client.EXPECT().GetTokenByPasswordGrant(gomock.Eq("testuser@test.local"), gomock.Eq("testpassword")).
		Times(1).Return(token, nil)

	Login(client, getLWSetting(), "testuser", "testpassword")
}

func TestEnsureUserProjects(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	token := oidc.NewMockTokenResponse(controller)
	token.EXPECT().GetAccessToken().Times(1).Return("some_access_token_string")
	u := models.User{Username: "testuser"}

	ensureUserProjects(getLWSetting(), u, token)
}

func TestGetGroups(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	verifier := tokenverifier.NewMockTokenVerifier(controller)
	token := oidc.NewMockTokenResponse(controller)
	token.EXPECT().GetAccessToken().Times(1).Return("some_access_token_string")

	jwtToken := &tokenverifier.JWTToken{
		Groups: []string{
			"test.Local\\\\Harbor_ProjectTestABC_Admin",
			"test.Local\\\\irrelevant_groupA",
			"test.local\\\\Users",
		}}
	verifier.EXPECT().Verify("some_access_token_string").Times(1).Return(jwtToken, nil)

	groupInfos := getUserProjectsFromToken(getLWSetting(), verifier, token)

	assert.NotNil(t, groupInfos)
	assert.Equal(t, len(groupInfos), 1)
	assert.Equal(t, groupInfos[0], GroupRole{groupName: "TestABC", role: "Admin"})
}
