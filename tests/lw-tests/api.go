package lwtest

import(
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"

    "github.com/vmware/harbor/src/common/models"
)

func getJSON(baseURL string, endpoint string, cookie http.Cookie) *http.Response {
    u, _ := url.ParseRequestURI(baseURL)
    u.Path = endpoint
    urlStr := fmt.Sprintf("%v", u)
    client := &http.Client{}

    r, _ := http.NewRequest("GET", urlStr, nil)
    r.Header.Add("Content-Type", "application/json")
    r.AddCookie(&cookie)

    resp, _ := client.Do(r)
    return resp
}

func postJSON(baseURL string, endpoint string, jsonData []byte, cookie http.Cookie) *http.Response {
    u, _ := url.ParseRequestURI(baseURL)
    u.Path = endpoint
    urlStr := fmt.Sprintf("%v", u)

    client := &http.Client{}

    r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonData))
    r.Header.Add("Content-Type", "application/json")
    r.AddCookie(&cookie)

    resp, _ := client.Do(r)
    return resp
}

func getResult(baseURL, endpoint, sessionID string, queryParams map[string]string) ([]byte, error){
    cookie := http.Cookie{Name: "beegosessionID", Value: sessionID}

    u, _ := url.ParseRequestURI(baseURL)
    u.Path = endpoint
    urlStr := fmt.Sprintf("%v", u)
    client := &http.Client{}

    r, _ := http.NewRequest("GET", urlStr, nil)

    query := r.URL.Query()
    for key, value := range queryParams {
            query.Add(key, value)
        }
    r.URL.RawQuery = query.Encode()

    r.AddCookie(&cookie)

    resp, _ := client.Do(r)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       panic(err.Error())
    }
    defer resp.Body.Close()
    return body, err
}

// Login uses the username and password combination to log in the user and
// returns the Session ID
func Login(baseURL, username, password string) (resp *http.Response, err error) {
    resource := "/login"
    data := url.Values{}
    data.Set("principal", username)
    data.Add("password", password)

    u, _ := url.ParseRequestURI(baseURL)
    u.Path = resource
    urlStr := fmt.Sprintf("%v", u)
    client := &http.Client{}
    r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    resp, err = client.Do(r)
    return
}

// CreateUser is used to create a User using the REST API
func CreateUser(baseURL string, sessionID string, user models.User) (bool, int) {
    var userCreated = false
    postData, _ := json.Marshal(user)
    cookie := http.Cookie{Name: "beegosessionID", Value: sessionID}
    resp := postJSON(baseURL, "/api/users", postData, cookie)

    if resp.StatusCode == 201 {
        userCreated = true
        fmt.Printf("User %s created successfully\n", user.Username)
    }
    defer resp.Body.Close()
    return userCreated, resp.StatusCode
}

// CreateProject is used to create a Project for a user using the REST API
func CreateProject(baseURL string, sessionID string, project ProjectReq) (bool, int) {
    var projectCreated = false
    cookie := http.Cookie{Name: "beegosessionID", Value: sessionID}
    jsonData, err := json.Marshal(project)
    if err != nil {
        panic(err)
    }

    resp := postJSON(baseURL, "/api/projects", jsonData, cookie)
    if resp.StatusCode == 201 {
        projectCreated = true
    } else if resp.StatusCode == 409 {
        // project name already exists
    }
    defer resp.Body.Close()
    return projectCreated, resp.StatusCode
}

// ListUsers fetches the users created in Harbor
func ListUsers(baseURL, sessionID string) []models.User {
    cookie := http.Cookie{Name: "beegosessionID", Value: sessionID}
    resp := getJSON(baseURL, "/api/users", cookie)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       panic(err.Error())
    }
    defer resp.Body.Close()
    var users = new([]models.User)
    json.Unmarshal([]byte(body), &users)
    return *users
}

// ListProjects fetches the projects created in Harbor
func ListProjects(baseURL, sessionID string) []models.Project {
    var queryParams = map[string]string{
        "page"      : "1",
        "page_size" : "100",
    }

    body, err := getResult(baseURL, "/api/projects", sessionID, queryParams)
    if err != nil {
       panic(err.Error())
    }
    var projects = new([]models.Project)
    json.Unmarshal([]byte(body), &projects)
    return *projects
}

func ListRepos(baseURL, sessionID string, projectID int64) []string {
    var queryParams = map[string]string{
        "project_id": strconv.FormatInt(projectID, 10),
        "page"      : "1",
        "page_size" : "100",
    }

    body, err := getResult(baseURL, "/api/repositories", sessionID, queryParams)
    if err != nil {
       panic(err.Error())
    }
    var repos = new([]string)
    json.Unmarshal([]byte(body), &repos)
    return *repos
}

// GetCurrentUser returns the currently logged in user
func GetCurrentUser(baseURL, sessionID string) models.User {
    cookie := http.Cookie{Name: "beegosessionID", Value: sessionID}
    resp := getJSON(baseURL, "/api/users/current", cookie)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       panic(err.Error())
    }
    defer resp.Body.Close()
    var currentUser = new(models.User)
    json.Unmarshal([]byte(body), &currentUser)
    return *currentUser
}

// SearchForProjects fetches a list of projects matching the search key supplied
func SearchForProjects(baseURL, sessionID string, searchKey string) []Search {
    var queryParams = map[string]string{
        "q": searchKey,
    }

    body, err := getResult(baseURL, "/api/search", sessionID, queryParams)
    if err != nil {
       panic(err.Error())
    }
    var searchedProjects = new([]Search)
    json.Unmarshal([]byte(body), &searchedProjects)
    return *searchedProjects
}
