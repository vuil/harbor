package lwtest

import (
        "math/rand"
        "time"

        "github.com/vmware/harbor/src/common/models"
        "github.com/manveru/faker"
)

// DataGen contains the test data to be persisted to Harbor
type DataGen struct {
    Users []models.User
    UserProjectMapping map[string][]ProjectImage
}

// ProjectImage contains the ProjectReq and images to be pushed
type ProjectImage struct {
    ProjectReq ProjectReq
    ImageAndVersions []ImageVersion
}

// ImageVersion contains the image name and tag combination to be pushed
type ImageVersion struct {
    ImageName string
    TagName string
}

// Globals
var images = []string{
    "hello-world",
    "alpine",
    "busybox"}

var tags = []string{
    "experimental",
    "test",
    "buggy",
    "1.1",
    "2.1",
}

// NewDataGen initializes an instance of DataGen
func NewDataGen() *DataGen {
    data := new(DataGen)
    data.UserProjectMapping = make(map[string][]ProjectImage)
    //data.Users = make([]models.User, noOfUsers)
    return data
}

// GetAllProjects returns the ProjectReqs while data generation
func (d *DataGen) GetAllProjects() []ProjectImage {
    var projectImages []ProjectImage
    for _, userProjects := range d.UserProjectMapping {
        for _, projImage := range userProjects {
            projectImages = append(projectImages, projImage)
        }
    }
    return projectImages
}

// GetProjects returns the set of Projects associated with that username
func (d *DataGen) GetProjects(username string) []ProjectImage {
    return d.UserProjectMapping[username]
}

func (d *DataGen) addProject2User(username string, project ProjectImage) {
    projects, _ := d.UserProjectMapping[username]
    projects = append(projects, project)
    d.UserProjectMapping[username] = projects
}

// This instantiates a Faker module used to generate randomized data
func getInstance() *faker.Faker {
    fake, err := faker.New("en")
    if err != nil {
      panic(err)
    }
    return fake
}

func _randInt(min, max int) int {
    return min + rand.Intn(max-min)
}

func _genProjectName() string {
    faker := getInstance()
    var name = faker.Characters(_randInt(5, 20))
    return name
}

func generateProjectReq(isPublic int32) ProjectReq {
    projectName := _genProjectName()
    return ProjectReq{
        ProjectName: projectName,
        Public: isPublic,
    }
}

func generateImageVersion() ImageVersion {
    imageIndex := _randInt(0, len(images))
    tagIndex := _randInt(0, len(tags))
    return ImageVersion {
        ImageName   : images[imageIndex],
        TagName     : tags[tagIndex],
    }
}

func generateProjectImage(isPublic int32, noOfImages int) ProjectImage {
    projectReq := generateProjectReq(isPublic)
    var imageVersions = make([]ImageVersion, noOfImages)
    for i := 0; i < noOfImages; i++ {
        imageVersions[i] = generateImageVersion()
    }
    return ProjectImage {
        ProjectReq      : projectReq,
        ImageAndVersions: imageVersions,
    }
}

// GenerateData generates test data to be persisted to Harbor
// Generates projects and images against Lightwave users
func (d *DataGen) GenerateProjectData(users []LWUserInfo, numberOfProjects, numberOfImages int) {
    for _, user := range users {
        // Seeding the random number generator
        rand.Seed(time.Now().UTC().UnixNano())
        for i := 0; i < numberOfProjects; i++ {
            // TODO: projects can be public or non-public
            var project = generateProjectImage(1, numberOfImages)
            d.addProject2User(user.Username, project)
        }
    }
}
