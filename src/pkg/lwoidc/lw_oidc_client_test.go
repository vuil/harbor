package lwoidc

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"

	"github.com/vmware/harbor/src/pkg/oidc"
)

func initClientWithMockup() oidc.Client {
	oidcOptions := oidc.ConfigOptions{
		DomainName:         "test.local",
		Endpoint:           "https://test.lw.com",
		AdminUser:          "AdminUser",
		AdminPassword:      "AdminPassword",
		IgnoreCertificates: false,
		TokenScopes:        "scopes",
		RootCAs:            nil,
	}

	factory := NewClientFactory(oidcOptions)
	client := factory.GetClient()
	lwclient := client.(*Client)
	httpClient := lwclient.GetHTTPClient()
	gock.InterceptClient(httpClient)

	return client
}

func TestGetTokenByPassword(t *testing.T) {
	defer gock.Off()
	client := initClientWithMockup()

	gock.New("https://test.lw.com").
		Post("openidconnect/token").
		Reply(200).
		JSON(map[string]interface{}{
			"access_token":  "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJBZG1pbmlzdHJhdG9yQHZzcGhlcmUubG9jYWwiLCJhdWQiOlsiQWRtaW5pc3RyYXRvckB2c3BoZXJlLmxvY2FsIiwicnNfYWRtaW5fc2VydmVyIl0sInNjb3BlIjoiYXRfZ3JvdXBzIHJzX2FkbWluX3NlcnZlciBvcGVuaWQgb2ZmbGluZV9hY2Nlc3MgaWRfZ3JvdXBzIiwiaXNzIjoiaHR0cHM6XC9cLzEwLjIwLjcyLjE5Mlwvb3BlbmlkY29ubmVjdFwvdnNwaGVyZS5sb2NhbCIsImdyb3VwcyI6WyJ2c3BoZXJlLmxvY2FsXFxVc2VycyIsInZzcGhlcmUubG9jYWxcXEFkbWluaXN0cmF0b3JzIiwidnNwaGVyZS5sb2NhbFxcQ0FBZG1pbnMiLCJ2c3BoZXJlLmxvY2FsXFxFdmVyeW9uZSJdLCJ0b2tlbl9jbGFzcyI6ImFjY2Vzc190b2tlbiIsInRva2VuX3R5cGUiOiJCZWFyZXIiLCJleHAiOjE0ODg5ODQyNTUsImlhdCI6MTQ4ODk4Mzk1NSwianRpIjoia19UZndmekI5V1ZObTB1bHJQcXl4MUdhYXVDYnVNdGNTTTBEYlBnaGkwcyIsInRlbmFudCI6InZzcGhlcmUubG9jYWwiLCJhZG1pbl9zZXJ2ZXJfcm9sZSI6IkFkbWluaXN0cmF0b3IifQ.fZd9evyN54dfwFfzf8sLQxmWqceJsk1ebJlDkfkgpjeI3-nA1fNTw9hHgfiLWTOXUcO_uKezVZd8ZOb7GsB5A_TK-GZ7kyEktXjCqdA1sSxVWa2Tlq4ICaR6MOeCG9zJYBtJYFD9q6zoH56Gk0t2phYgNilTjjINBLfyBIGJnIwx0aalm6Lm3EOuTGxQyMSHVHDMRL_En8phpgZUfiSBfEw9f6dpWE2xWPFQJ_TB8N1lRsVnIQAnf0OG7zw3y4_pxMNQ5L11b-vhIyQzxh6HzZ8m6fcfP8nPehpT6B9tpNu7czjTpBcNYvOCXWQIc00IIGig2jN17ZXE6W-JTHlpaDvtzjnt96XsOny4yF6pg3FdRc9tSDLH-A_lkWDTEOYF7633UWJK7UC-WQE9k9ecZhjMAX3UsgriuzuviWmK1dxG0JICzwXPAp1S6FdytGTjam113_MpTnF4Dr91yRM2s78_FRoc88DF6qmgFiiuO_6EsaIwDIfy6o3apE5j_G7JR-lkohJAzguOoyFxPL_2fhDbZOQRhaYHBDjhXS7hRiNi8nSxF711KP123vzFeRd6riDCIRBy0onLJoe2xivxO-ee_XlkzrlHdkZzowF1pPSadrK9rL7dTREboeH1F8b-9qsbA3Zil2s2yu13n_TO_J_j7NImllTBcBO8ama3exE",
			"refresh_token": "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJBZG1pbmlzdHJhdG9yQHZzcGhlcmUubG9jYWwiLCJhdWQiOiJBZG1pbmlzdHJhdG9yQHZzcGhlcmUubG9jYWwiLCJzY29wZSI6ImF0X2dyb3VwcyByc19hZG1pbl9zZXJ2ZXIgb3BlbmlkIG9mZmxpbmVfYWNjZXNzIGlkX2dyb3VwcyIsImlzcyI6Imh0dHBzOlwvXC8xMC4yMC43Mi4xOTJcL29wZW5pZGNvbm5lY3RcL3ZzcGhlcmUubG9jYWwiLCJ0b2tlbl9jbGFzcyI6InJlZnJlc2hfdG9rZW4iLCJ0b2tlbl90eXBlIjoiQmVhcmVyIiwiZXhwIjoxNDg5MDA1NTU1LCJpYXQiOjE0ODg5ODM5NTUsImp0aSI6IkRCZUEzcG1wWVdxa0Z2ZklXbE1FMXdNTWQwQjg2Snlvb1N6YlV2dFZQRm8iLCJ0ZW5hbnQiOiJ2c3BoZXJlLmxvY2FsIn0.DTn1x7-RlsfN8NRITRoe2eU2p8BfdmNJDECPtUyVaREGEug1zarRvd4jCoTiXUQQeQWXjhatnod2OCZJ37xO17prX2Pe_mOxAP-3rLAZL3a0Niuh53NdoajcDd2qwoObMbTclAiIQw-M9Wv0jH5-IThFHesaHmmxjRoafYBTHTvEVuumAQCksfETZoL0qBOlY7oNkvS3O817utOgZpbWUfGg_Ts3ZJGmBDCG5YGVEiQUcxm9-2wDs1nOsYjYgsdLovw53-P1AaPoJXJ7zbGoSVn_4gnYZh2zlsfI4kPBxV4hoo5bCpkrOVb8TIN7O3c40I9Rwe3ph8EFh79qlyTuxr7jUW_v3U0bvr9AXAtsXHUMRiI42IBb5M7EPA4pjPvXqI5yF8BzKv26xdrKdd1HeAVCNIZSnG8qNQN9sec6zcP0Uu0ptJPpTFklOtgUgK1_K94JfPTr1ZVsB_2ZpwXigfxCzTND60Mq_1i623h2nnZ6klvjDPeHZhxHHC5AzP7DtNIVuZxA--CbkMsdoUJEkMyTHE4O8wvzV62M5WOk_Gaeu77t_7O3DrhCOnAV1yGcf81OVPMRNjsFAcUrYr2N5f4IjvsVouxPWnG-o6Gp5KCvXGuAb7Gg_Fcp8vJBawopCECKq8a7VROsWQcOXE3HIx6I0XhYRpU2dyHbYVB9d7g",
			"id_token":      "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJBZG1pbmlzdHJhdG9yQHZzcGhlcmUubG9jYWwiLCJpc3MiOiJodHRwczpcL1wvMTAuMjAuNzIuMTkyXC9vcGVuaWRjb25uZWN0XC92c3BoZXJlLmxvY2FsIiwiZ3JvdXBzIjpbInZzcGhlcmUubG9jYWxcXFVzZXJzIiwidnNwaGVyZS5sb2NhbFxcQWRtaW5pc3RyYXRvcnMiLCJ2c3BoZXJlLmxvY2FsXFxDQUFkbWlucyIsInZzcGhlcmUubG9jYWxcXEV2ZXJ5b25lIl0sInRva2VuX2NsYXNzIjoiaWRfdG9rZW4iLCJ0b2tlbl90eXBlIjoiQmVhcmVyIiwiZ2l2ZW5fbmFtZSI6IkFkbWluaXN0cmF0b3IiLCJhdWQiOiJBZG1pbmlzdHJhdG9yQHZzcGhlcmUubG9jYWwiLCJzY29wZSI6ImF0X2dyb3VwcyByc19hZG1pbl9zZXJ2ZXIgb3BlbmlkIG9mZmxpbmVfYWNjZXNzIGlkX2dyb3VwcyIsImV4cCI6MTQ4ODk4NDI1NSwiaWF0IjoxNDg4OTgzOTU1LCJmYW1pbHlfbmFtZSI6InZzcGhlcmUubG9jYWwiLCJqdGkiOiIwUC1BVHRwSzB5dTB4UzZZUHBJeVZYWkVsQ18yVjlvZlo1Rk5rTmR6dk13IiwidGVuYW50IjoidnNwaGVyZS5sb2NhbCJ9.AeNKluQm_WdwBftGUrj8xioVSe0kcwiDFxUSqn2cXbTs1vRvyKEqtWgT46gaFQ3Fy4B4e659rFHdjuKR9jUlOkRlqloCH67Ei8fsPzmoeIk7-kXSGz7KgVZZsPie588ECXpoBbpFZSaBgL07aoxJoXM2Fx1PyjUQUncxQcMcC_OCy7Z6obr13qK9nFuTKPepC0IiFQ6NxruKbjh_k7Uc9NNKWyAbeMFr5bM-Q53kp1lxzunfs5JUvIqRww_12CvVDKJ8flRMc2qL1L9YpEqzwXmYw9G9Uvi-wwjpeVOpQvwIY43pAx7GmvDKzzlBMWITAr-TOUftNqiLjtuY-i150_ZAmTUhisOK8YRN9EwRNC0DmQTt5unuPVnF7q5v5GDfCc53qu8zLTxsTjCxqajXnnRtIelNgCDsaGQEdxh3UOqWR3tPVOA-ndm6ELzv7SDCAm3K_5QQbZg2J5jDqExoU-AXYjX9P6LsDXMDz1u3qug7IyQDyhzFP5OHC4HMLbrl4bsDL4IxTXzrnnsGPccgdRUOUWpPJwzVC7Ibgk13G7S0NaAsCeLAY0oEcBXJdVl15he6PFpj66xpS8QAwa8x3R6-h9jmpb7njpkDYX-05ZCcimMkA2dvSElSr-E73-zeOylCgHXkQpyzZf2zU1ut4pl0W829FEZbOSLbMfLQ5xU",
			"token_type":    "Bearer",
			"expires_in":    300})
	response, errorResponse := client.GetTokenByPasswordGrant("testuser", "testpassword")
	assert.Nil(t, errorResponse)
	assert.NotNil(t, response)
}

func TestCreateGroupAlreadyExists(t *testing.T) {
	defer gock.Off()
	client := initClientWithMockup()
	token := &TokenResponse{}

	gock.New("https://test.lw.com").
		Post("vmdir/tenant/test.local/groups").
		Reply(400).
		JSON(map[string]string{
			"error":   "bad_request",
			"details": "Failed to create group 'Harbor_ProjectX_Dev' on tenant 'vsphere.local'",
			"cause":   "Another user or group Harbor_ProjectX_Dev already exists with the same name"})

	response := client.CreateGroup("Harbor_ProjectX_Dev", "project x dev role", token)
	assert.NotNil(t, response.GetErrorResponse())
}

func TestCreateroup(t *testing.T) {
	defer gock.Off()
	client := initClientWithMockup()
	token := &TokenResponse{}

	gock.New("https://test.lw.com").
		Put("vmdir/tenant/test.local/groups/Harbor_ProjectX_Dev@test.local/members").
		Reply(200).
		JSON(map[string]interface{}{
			"name":    "Harbor_ProjectXYZ4_Dev",
			"domain":  "vsphere.local",
			"details": map[string]interface{}{"description": "project XYZ dev role"}})

	response := client.CreateGroup("Harbor_ProjectX_Dev", "project x dev role", token)
	assert.Nil(t, response.GetErrorResponse())
}

func TestAddUserToGroupAlreadyExists(t *testing.T) {
	defer gock.Off()
	client := initClientWithMockup()
	token := &TokenResponse{}

	gock.New("https://test.lw.com").
		Put("vmdir/tenant/test.local/groups/Harbor_ProjectX_Dev@test.local/members").
		MatchParams(map[string]string{"type": "user", "members": "tester@test.local"}).
		Reply(400).
		JSON(map[string]string{
			"error":   "bad_request",
			"details": "Failed to add groups to group 'Harbor_ProjectX_Dev@vsphere.local' in tenant 'test.local'",
			"cause":   "group Harbor_ProjectX_Dev currently has user CN=tester acct,cn=users,dc=test,dc=local as its member"})

	response := client.AddUserToGroup("Harbor_ProjectX_Dev", "tester", token)
	assert.NotNil(t, response.GetErrorResponse())
}

func TestAddUserToGroup(t *testing.T) {
	defer gock.Off()
	client := initClientWithMockup()
	token := &TokenResponse{}

	gock.New("https://test.lw.com").
		Put("vmdir/tenant/test.local/groups/Harbor_ProjectX_Dev@test.local/members").
		MatchParams(map[string]string{"type": "user", "members": "tester@test.local"}).
		Reply(200).
		JSON(map[string]string{})

	response := client.AddUserToGroup("Harbor_ProjectX_Dev", "tester", token)
	assert.Nil(t, response.GetErrorResponse())
}
