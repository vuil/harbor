package main

import(
    "fmt"
    "github.com/vmware/harbor/src/common/models"
)

func CompareUsers(u1, u2 models.User) bool {
    if &u1 == nil && &u2 == nil { return true }
    if &u1 == nil || &u2 == nil { return false }

    isEqual := u1.Username == u2.Username
    if isEqual { isEqual = u1.Email == u2.Email }
    if isEqual { isEqual = u1.Realname == u2.Realname }
    if isEqual { isEqual = u1.Comment == u2.Comment }

    return isEqual
}

func CompareProjects(p1 models.Project, p2 ProjectReq, user models.User) bool {
    if &p1 == nil { return false }

    isEqual := p1.Name == p2.ProjectName
    if isEqual { isEqual = int32(p1.Public) == p2.Public }
    if isEqual { isEqual = p1.OwnerID == user.UserID }

    return isEqual
}

func CompareRepos(projectName string, imgVersion ImageVersion, repositories []string) bool {
    var isEqual = false
    repoName := fmt.Sprintf("%s/%s", projectName, imgVersion.ImageName)
    for _, repo := range repositories {
        isEqual = repoName == repo
        if isEqual {
            break;
        }
    }
    return isEqual
}
