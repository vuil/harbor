package lwtest

import(
    "fmt"
    "strconv"
    "strings"

    "github.com/codeskyblue/go-sh"
)

// LwHarborInstance contains the parameters for a LW integrated Harbor instance
type LwHarborInstance struct {
    BaseURL string
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
func (self *LwHarborInstance) Login(user LWUserInfo) (string, bool) {
    resp, err := Login(self.BaseURL, user.Username, user.Password)
    statusCode := resp.StatusCode
    if err != nil {
        fmt.Errorf("Error while logging in with credentials %s/%s", user.Username, user.Password)
        return "", false
    } else {
        if statusCode != 200 {
            fmt.Errorf("Error while logging in with credentials %s/%s", user.Username, user.Password)
            return "", false
        } else {
            sessionID := resp.Cookies()[0].Value
            return sessionID, true
        }
    }
}

// CreateProject creates a Project entry in the Harbor db
func (self *LwHarborInstance) CreateProject(user LWUserInfo, project ProjectReq) bool {
    sessionID, isLoggedIn := self.Login(user)
    if isLoggedIn {
        created, respCode := CreateProject(self.BaseURL, sessionID, project)
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
func (self *LwHarborInstance) PushImage(user LWUserInfo, project, image, tag string) bool {
    endpoint := parseHarborURL(self.BaseURL)
    return pushImage(user.Username, user.Password, endpoint, strconv.Itoa(self.RegistryPort), project, image, tag)
}

func pushImage(username, password, endpoint, port, project, image, tag string) bool {
    session := sh.NewSession()
    session.ShowCMD = true
    err := session.Command("docker", "pull", image).Run()
    if err == nil { err = session.Command("docker", "login", "-u", username, "-p", password, ""+endpoint+":"+port).Run() }
    if err == nil { err = session.Command("docker", "tag", image, ""+endpoint+":"+port+"/"+project+"/"+image+":"+tag).Run() }
    if err == nil { err = session.Command("docker", "push", ""+endpoint+":"+port+"/"+project+"/"+image+":"+tag).Run() }
    if err != nil { return false }
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
