package lwtest

import(
    "flag"
    "os"
    "strconv"
    "testing"
)

// PersistenceTestParams
type PersistenceTestParams struct {
    NumberOfProjects int
    NumberOfImages int
}

// LWUserInfo contains the initial username/pwd combos in the test LW instance
type LWUserInfo struct {
    Username string
    Password string
}

var (
    testingParams = map[string]PersistenceTestParams{
        "short" : {5, 5},
        "medium": {10, 20},
        "large" : {25, 50},
    }

    lwUsers = []LWUserInfo {
        {Username: "harbor.test1", Password: "VMware@123"},
        {Username: "harbor.test2", Password: "VMware@123"},
    }

    data *DataGen
)

// LwSetup contains the user and repo combo for testing initial LW setup
type LwSetup struct {
    User LWUserInfo
    RepoName string
}

// Pre-populated Lightwave users and projects
/*var LWUserRepos = []LwSetup{
    {UserName: "harbor.test1", RepoName: "HarborTestOne"},
    {UserName: "harbor.test2", RepoName: "HarborTestTwo"},
}*/

func TestLoginForLWUsers(t *testing.T) {
    port, _ := strconv.Atoi(registryPort)
    instance := NewLwHarborInstance(harborEndpoint, port)
    for _, user := range lwUsers {
        t.Logf("Logging in with default LW user %s\n", user.Username)
        _, login := instance.Login(user)
        if !login {
            t.Errorf("Login failed")
        }
    }
}

func TestCreationOfProjects(t *testing.T) {
    port, _ := strconv.Atoi(registryPort)
    instance := NewLwHarborInstance(harborEndpoint, port)
    for username, projectImages := range data.UserProjectMapping {
        var user LWUserInfo
        for _, ui := range lwUsers { if ui.Username == username { user = ui } }
        for _, projectImage := range projectImages {
            t.Logf("Creating project %s with user %s\n", projectImage.ProjectReq.ProjectName, username)
            isCreated := instance.CreateProject(user, projectImage.ProjectReq)
            if !isCreated {
                t.Errorf("Project creation failed\n")
            }
        }
    }
}

func TestPushingImages(t *testing.T) {
    port, _ := strconv.Atoi(registryPort)
    instance := NewLwHarborInstance(harborEndpoint, port)
    for username, projectImages := range data.UserProjectMapping {
        var user LWUserInfo
        for _, ui := range lwUsers { if ui.Username == username { user = ui } }
        for _, projectImage := range projectImages {
            var projectName = projectImage.ProjectReq.ProjectName
            for _, imageAndVersion := range projectImage.ImageAndVersions {
                imgPushed := instance.PushImage(user, projectName, imageAndVersion.ImageName, imageAndVersion.TagName)
                t.Logf("Pushing image %s with tag %s to project %s using user %s\n", imageAndVersion.ImageName, imageAndVersion.TagName, projectName, user.Username)
                if !imgPushed {
                    t.Errorf("Could not push image\n")
                }
            }
        }
    }
}

var (
    short = flag.Bool("short", true, "run short LW-Harbor persistence tests")
    medium  = flag.Bool("medium", false, "run medium LW-Harbor persistence tests")
    large = flag.Bool("large", false, "run larger LW-Harbor persistence tests")

    harborEndpoint = os.Getenv("HARBOR_INSTANCE")
    registryPort = os.Getenv("REGISTRY_PORT")

    pTestParams PersistenceTestParams
)

func TestMain(m *testing.M) {
    flag.Parse()
    // Only one flag should be passed, ignore multiple flags
    if *large {
        pTestParams = testingParams["large"]
    } else if *medium {
        pTestParams = testingParams["medium"]
    } else {
        pTestParams = testingParams["short"]
    }

    // Generate Harbor test data
    initiateDataGeneration(lwUsers, pTestParams)

    // Run the test suite
    result := m.Run()
    os.Exit(result)
}

func initiateDataGeneration(users []LWUserInfo, params PersistenceTestParams) *DataGen {
    data = NewDataGen()
    data.GenerateProjectData(users, params.NumberOfProjects, params.NumberOfImages)
    return data
}
