
package main

import (
        "flag"
        "fmt"
        "math/rand"
        "os"
        "os/exec"
        "path"
        "strings"
        "time"

        "github.com/vmware/harbor/src/common/models"
)

// Harbor deployment URL
var baseURL string

// This is to get the session ID using the admin login credentials
// Used to create the users in the system
func initialLogin() string {
    return Login(baseURL, "admin", "Harbor12345")
}

func createProjectForUser(user models.User, project ProjectReq) (bool, int) {
    sessionID := Login(baseURL, user.Username, user.Password)
    projectCreated, status := CreateProject(baseURL, sessionID, project)

    if projectCreated {
        fmt.Printf("Project named %s created successfully\n", project.ProjectName)
    }
    return projectCreated, status
}

func CToGoString(c []byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }
        n = i
    }
    return string(c[:n+1])
}

func pushImageToRepo(username, password, endpoint, port, project, image, tag string) bool {
    cwd, _ := os.Getwd()
    pushImagePath := path.Join(cwd, "push_image.sh")
    cmd := exec.Command(pushImagePath, username, password, endpoint, port, project, image, tag)
    output, err := cmd.CombinedOutput()
    if err != nil {
        //panic(err)
        s := CToGoString(output[:])
        fmt.Println(s)
        return false
    }
    return true
}

func parseHarborURL(url string) string {
    var parsedURL = url
    if strings.Contains(url, "https://") {
        parsedURL = strings.TrimLeft(baseURL, "https://")
    } else if strings.Contains(url, "http://") {
        parsedURL = strings.TrimLeft(baseURL, "http://")
    }
    return parsedURL
}

func PopulateData(data DataGen, baseURL string, portNumber string) {
    sessionID := initialLogin()

    for i := 0; i < len(data.Users); i++ {
        var userCreated = false
        for !userCreated {
            userCreated, _ = CreateUser(baseURL, sessionID, data.Users[i])
        }
    }

    for username, projectImages := range data.UserProjectMapping {
        user := data.GetUser(username)
        for _, projectImage := range projectImages {
            var projectCreated = false
            var status = 0
            for !projectCreated && status != 409 {
                projectCreated, status = createProjectForUser(user, projectImage.ProjectReq)
                if status == 409 {
                    fmt.Printf("Project named %s already exists\n", projectImage.ProjectReq.ProjectName)
                } else {
                    var projectName = projectImage.ProjectReq.ProjectName
                    for _, imageAndVersion := range projectImage.ImageAndVersions {
                        url := parseHarborURL(baseURL)
                        imgPushed := pushImageToRepo(user.Username, user.Password, url, portNumber, projectName, imageAndVersion.ImageName, imageAndVersion.TagName)
                        if imgPushed {
                            fmt.Printf("Pushed image %s with tag %s to project %s using user %s\n", imageAndVersion.ImageName, imageAndVersion.TagName, projectName, user.Username)
                        }
                    }
                }
            }
        }
    }
}

func exitWithMessage(message string) {
    fmt.Println(message)
    os.Exit(1)
}

func verifyProjectsForUser(user models.User, projects []models.Project, dataGen DataGen) (bool, string) {
    var projectImages = dataGen.GetProjects(user.Username)
    var isEqual = true
    var msg = ""
    for i:=0; i<len(projectImages) && isEqual; i++ {
        projReq := projectImages[i].ProjectReq
        //TODO: Remove the compared project from the projects[]
        for j:=0; j<len(projects); j++ {
            if projReq.ProjectName == projects[j].Name {
                isEqual = CompareProjects(projects[j], projReq, user)
                if isEqual {
                    isEqual, msg = verifyReposForProject(projects[j], projectImages[i])
                    if !isEqual {
                        return isEqual, msg
                    }
                    break;
                } else {
                    msg := fmt.Sprintf("Project named %s is not the same\n", projReq.ProjectName)
                    return isEqual, msg
                }
            }
        }
    }
    return isEqual, ""
}

func verifyReposForProject(project models.Project, projectImage ProjectImage) (bool, string){
    var sessionID = initialLogin()
    var repositories = ListRepos(baseURL, sessionID, project.ProjectID)

    isEqual := len(projectImage.ImageAndVersions) == len(repositories)
    for _, imgVer := range projectImage.ImageAndVersions {
        projectName := projectImage.ProjectReq.ProjectName
        isEqual = CompareRepos(projectName, imgVer, repositories)
        if isEqual {
            break;
        } else {
            msg := fmt.Sprintf("Repo named %s for project named %s not found", imgVer.ImageName, projectName)
            return isEqual, msg
        }
    }
    return isEqual, ""
}

func VerifyData(data DataGen) {
    sessionID := initialLogin()
    users := ListUsers(baseURL, sessionID)

    if len(users) != len(data.Users) {
        msg := fmt.Sprintf("Number of users fetched should be equal. Expected %d Actual %d\n", len(data.Users), len(users))
        exitWithMessage(msg)
    }
    projects := ListProjects(baseURL, sessionID)
    generatedProjects := data.GetAllProjects()
    // Adding 1 for the default project library
    if len(projects) != len(generatedProjects)+1 {
        msg := fmt.Sprintf("Number of projects fetched should be equal. Expected %d Actual %d\n", len(generatedProjects)+1, len(projects))
        exitWithMessage(msg)
    }

    for i := 0; i < len(data.Users); i++ {
        var username = users[i].Username
        var usr = data.GetUser(username)
        if !CompareUsers(users[i], usr) {
            msg := fmt.Sprintf("User info different")
            exitWithMessage(msg)
        } else {
            eq, msg := verifyProjectsForUser(users[i], projects, data)
            if !eq {
                exitWithMessage(msg)
            }
        }
    }
}

func VerifyUserPushPermissions(baseURL, portNumber string, data DataGen) {
    sessionID := initialLogin()
    users := ListUsers(baseURL, sessionID)
    projects := ListProjects(baseURL, sessionID)

    rand.Seed(time.Now().UTC().UnixNano())
    projectIndex := rand.Intn(len(projects))
    userIndex := rand.Intn(len(users))
    for projects[projectIndex].OwnerID == users[userIndex].UserID {
        userIndex = rand.Intn(len(users))
    }
    user := users[userIndex]
    imgVersion := data.GetImageVersion()

    url := parseHarborURL(baseURL)
    defaultPwd := "Harbor12345"
    fmt.Println("Trying to push image without adequate permissions")
    imgPushed := pushImageToRepo(user.Username, defaultPwd, url, portNumber, projects[projectIndex].Name, imgVersion.ImageName, imgVersion.TagName)
    if !imgPushed {
        fmt.Printf("Image push failed.\nApplying adequate permissions for the user\n")
        memberReq := MemberReq {
            Username: user.Username,
            Roles   : []int{ 2 }, // developer role
        }
        // TODO: Check response code for permissions already set (409)
        roleAdded, _ := AddPermissionsForUser(sessionID, projects[projectIndex].ProjectID, memberReq)
        if roleAdded {
            if pushImageToRepo(user.Username, defaultPwd, url, portNumber, projects[projectIndex].Name, imgVersion.ImageName, imgVersion.TagName) {
                fmt.Printf("Image push successful with user %s added as member to project %s.\n", user.Username, projects[projectIndex].Name)
            }
        } else {
            msg := fmt.Sprintf("Error occured while adding role to user %s.\n", user.Username)
            exitWithMessage(msg)
        }
    }
}

func main() {
    urlPtr := flag.String("harbor_url", "http://localhost", "Harbor REST endpoint")
    noOfUsersPtr := flag.Int("users", 10, "Number of users to be created")
    noOfProjectsPtr := flag.Int("projects", 100, "Number of projects to be created")
    portNumberPtr := flag.String("port", "5000", "Port number to access the endpoint, 80/5000")

    flag.Parse()
    baseURL = *urlPtr

    data := NewDataGen(*noOfUsersPtr)
    fmt.Println("Generating data...")
    data.GenerateData(*noOfUsersPtr, *noOfProjectsPtr)

    fmt.Println("Populating data...")
    PopulateData(*data, baseURL, *portNumberPtr)

    fmt.Println("Verifying populated data...")
    VerifyData(*data)

    fmt.Println("Verifying user permissions...")
    VerifyUserPushPermissions(baseURL, *portNumberPtr, *data)
}
