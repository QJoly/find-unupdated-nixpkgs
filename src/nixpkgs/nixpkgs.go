package nixpkgs

import (
	"fmt"
	"funixpkgs/src/github_api"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Package struct {
	pname    string
	owner    string
	repo     string
	version  string
	filepath string
}

func FindUnupdatedPkgs(path string) ([]string, error) {
	findNixPkgs, err := findNixPkgs(path)
	if err != nil {
		return nil, err
	}
	getPkgsData(findNixPkgs)
	return findNixPkgs, nil
}

func findNixPkgs(path string) ([]string, error) {
	var result []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "default.nix" && checkFetchFromGithub(filePath) {
			result = append(result, filePath)
		}
		return nil
	})

	return result, err
}

func getPkgsData(findNixPkgs []string) []Package {

	listPkgs := []Package{}
	for _, pkg := range findNixPkgs {

		Package := Package{
			pname:    "",
			version:  "",
			owner:    "",
			repo:     "",
			filepath: pkg,
		}

		fileContents, err := ioutil.ReadFile(pkg)
		if err != nil {
			panic(err)
		}

		fileString := string(fileContents)
		lines := strings.Split(fileString, "\n")

		for _, line := range lines {
			// Check if the line contains "pname"
			if strings.Contains(line, "pname") {

				if len(strings.Split(line, `"`)) > 1 {
					Package.pname = strings.Split(line, `"`)[1]
				}
			}
			if strings.Contains(line, "owner") {

				if len(strings.Split(line, `"`)) > 1 {
					Package.owner = strings.Split(line, `"`)[1]
				}
			}

			if strings.Contains(line, "version =") {

				if len(strings.Split(line, `"`)) > 1 {
					if checkIfVersionContainsPrefix(pkg) {
						Package.version = "v" + strings.Split(line, `"`)[1]
					} else {
						Package.version = strings.Split(line, `"`)[1]
					}
				}
			}
			if strings.Contains(line, "repo") {

				if len(strings.Split(line, `"`)) > 1 {
					Package.repo = strings.Split(line, `"`)[1]
				}
			}

		}
		listPkgs = append(listPkgs, Package)
	}
	for _, pkg := range listPkgs {
		//		fmt.Println(pkg)
		if github_api.CheckIfRepoExist(pkg.owner, pkg.pname) {
			github_version, err := github_api.GetLatestRelease(pkg.owner, pkg.pname)
			if err != nil || github_version == "0.0.0" {
				continue
			}
			if pkg.version != github_version {
				fmt.Printf("%s - %s\n", pkg.pname, pkg.filepath)
				fmt.Printf("Latest version: %s - NixPkgs version : %s\n", github_version, pkg.version)
				fmt.Println(strings.Repeat("*", 10))
			}

		}
	}

	return listPkgs
}

func checkFetchFromGithub(path string) bool {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fileString := string(fileContents)

	pattern := regexp.MustCompile(`fetchFromGitHub`)

	if pattern.MatchString(fileString) {
		return true
	} else {
		return false
	}
}

func checkIfVersionContainsPrefix(path string) bool {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	fileString := string(fileContents)

	pattern := regexp.MustCompile(`rev\s*=\s*"v\${version}";`)

	if pattern.MatchString(fileString) {
		return true
	} else {
		return false
	}
}
