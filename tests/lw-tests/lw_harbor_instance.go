package lwtest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codeskyblue/go-sh"

	"github.com/vmware/harbor/src/common/models"
)

// LwHarborInstance contains the parameters for a LW integrated Harbor instance
type LwHarborInstance struct {
	BaseURL      string
	RegistryPort int
}

// NewLwHarborInstance is the constructor for LwHarborInstance
func NewLwHarborInstance(baseURL string, port int) *LwHarborInstance {
	instance := new(LwHarborInstance)
	instance.BaseURL = baseURL
	instance.RegistryPort = port
	return instance
}

// Login method logs in the user
func (instance *LwHarborInstance) Login(user LWUserInfo) (string, bool) {
	resp, err := Login(instance.BaseURL, user.Username, user.Password)
	statusCode := resp.StatusCode
	if err != nil {
		fmt.Printf("Error while logging in with credentials %s/%s", user.Username, user.Password)
		return "", false
	} 
	if statusCode != 200 {
		fmt.Printf("Error while logging in with credentials %s/%s", user.Username, user.Password)
		return "", false
	} 
	sessionID := resp.Cookies()[0].Value
	return sessionID, true
	}
}

// CreateProject creates a Project entry in the Harbor db
func (instance *LwHarborInstance) CreateProject(user LWUserInfo, project ProjectReq) bool {
	sessionID, isLoggedIn := instance.Login(user)
	if isLoggedIn {
		created, respCode := CreateProject(instance.BaseURL, sessionID, project)
		if !created {
			if respCode == 409 {
				fmt.Printf("Project named %s already exists. Could not recreate project\n", project.ProjectName)
			} else {
				fmt.Errorf("Error while creating project %s", project.ProjectName)
			}
		}
		return created
	}
	return false
}

// PushImage pushes an image to the repository
func (instance *LwHarborInstance) PushImage(user LWUserInfo, project, image, tag string) bool {
	endpoint := parseHarborURL(instance.BaseURL)
	return pushImage(user.Username, user.Password, endpoint, strconv.Itoa(instance.RegistryPort), project, image, tag)
}

func (instance *LwHarborInstance) AddPermissions(user LWUserInfo, project string) bool {
	sessionID, loggedIn := instance.Login(LWUserInfo{"admin", "Harbor12345"})
	if !loggedIn {
		fmt.Println("Login with admin user failed")
		return false
	}
	projects := instance.GetProjects(project)
	// TODO: Handle error response codes
	addedPerms, _ := AddPermissionsForUser(instance.BaseURL, sessionID, projects[0].ProjectID, MemberReq{
		Username: user.Username,
		Roles:    []int{2}, // developer role
	})
	if !addedPerms {
		fmt.Println("Could not add permissions for user")
	}
	return addedPerms
}

func (instance *LwHarborInstance) GetProjects(projectName string) []models.Project {
	sessionID, loggedIn := instance.Login(LWUserInfo{"admin", "Harbor12345"})
	if !loggedIn {
		fmt.Println("Login with admin user failed")
		return []models.Project{}
	}
	return ListProjects(instance.BaseURL, sessionID, projectName)
}

func pushImage(username, password, endpoint, port, project, image, tag string) bool {
	session := sh.NewSession()
	session.ShowCMD = true
	err := session.Command("docker", "pull", image).Run()
	if err == nil {
		err = session.Command("docker", "login", "-u", username, "-p", password, ""+endpoint+":"+port).Run()
	}
	if err == nil {
		err = session.Command("docker", "tag", image, ""+endpoint+":"+port+"/"+project+"/"+image+":"+tag).Run()
	}
	if err == nil {
		err = session.Command("docker", "push", ""+endpoint+":"+port+"/"+project+"/"+image+":"+tag).Run()
	}
	if err != nil {
		return false
	}
	return true
}

func parseHarborURL(url string) string {
	var parsedURL = url
	if strings.Contains(url, "https://") {
		parsedURL = strings.TrimLeft(url, "https://")
	} else if strings.Contains(url, "http://") {
		parsedURL = strings.TrimLeft(url, "http://")
	}
	return parsedURL
}
