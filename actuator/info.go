package actuator

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
)

var gitInfo map[string]string = make(map[string]string)

func loadGitInfo() {
	file, err := os.Open("git.properties")


	if err != nil {
		// can't open file or not exists
		out, err := exec.Command("bash", "-c", "git rev-parse --abbrev-ref HEAD; git show -s --format=\"%h%n%ci\"").Output()
		if err != nil {
			log.Println(err)
		}
		splitted := strings.Split(string(out), "\n")
		if len(splitted) > 3 {
			gitInfo["branch"] = splitted[0]
			gitInfo["commitid"] = splitted[1]
			gitInfo["committime"] = splitted[2]
		}

	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		gitInfo["branch"] = scanner.Text()
		scanner.Scan()
		gitInfo["commitid"] = scanner.Text()
		scanner.Scan()
		gitInfo["committime"] = scanner.Text()
	}

}

func info() string {
	infoJson := generateInfoData()
	b, err := json.Marshal(infoJson)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func generateInfoData() *infoJson {
	infoJson := new(infoJson)
	infoJson.Git = new(git)
	infoJson.Git.Commit = new(commit)

	infoJson.Git.Branch = gitInfo["branch"]
	infoJson.Git.Commit.ID = gitInfo["commitid"]
	infoJson.Git.Commit.Time = gitInfo["committime"]

	return infoJson
}

type infoJson struct {
	Git *git `json:"git"`
}

type git struct {
	Branch string `json:"branch"`
	Commit *commit `json:"commit"`
}

type commit struct {
	ID string `json:"id"`
	Time string `json:"time"`
}