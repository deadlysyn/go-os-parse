package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

func TestMain(m *testing.M) {
	if setupSubtest() {
		os.Exit(222)
	}
	exitVal := m.Run()

	err := teardownSubtest()
	if err != nil {
		fmt.Printf("welp teardown went terribly...not sure if it matters or not but here is your error: %v", err)
		os.Exit(333)
	}

	os.Exit(exitVal)
}

func TestDockerIntegration(t *testing.T) {
	tests := map[string]struct {
		expectedDockerfile string
		expectedResult     string
	}{
		"alpine": {
			expectedDockerfile: "alpine/Dockerfile",
			expectedResult:     "apk",
		},
		"arch": {
			expectedDockerfile: "arch/Dockerfile",
			expectedResult:     "pacman",
		},
		"debian": {
			expectedDockerfile: "debian/Dockerfile",
			expectedResult:     "dpkg",
		},
		"fedora": {
			expectedDockerfile: "fedora/Dockerfile",
			expectedResult:     "dnf",
		},
		"centos": {
			expectedDockerfile: "centos/Dockerfile",
			expectedResult:     "yum",
		},
	}
	t.Run("docker", func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// faster, harder to read...
				t.Parallel()
				t.Logf(">>>>>>>>>>>>>> %v: Started", name)
				t.Logf(">>>>>>>>>>>>>> %v: %v", name, test.expectedDockerfile)
				output, status := runCommand("docker", "build", "--no-cache", "-f", test.expectedDockerfile, ".")
				t.Logf(">>>>>>>>>>>>>> %v: %v", name, output)
				if status == false && !strings.Contains(output, "package manager") {
					t.Errorf(">>>>>>>>>>>>>> %v: Docker build failed and was not due to detection. See test output for more details.", name)
					return
				}
				if !strings.Contains(output, test.expectedResult) {
					t.Errorf(">>>>>>>>>>>>>> %v: Docker build failed. See test output for details.", name)
				}
			})
		}
	})
}

func teardownSubtest() error {
	fmt.Println("[TEARDOWN]")
	return os.Remove("osdetect")
}

func setupSubtest() bool {
	fmt.Println("[SETUP]")
	output, goStatus := runCommand("go", "build", "-o", "osdetect", "../main.go")
	if goStatus == false {
		fmt.Println("Could not build osdetect")
		fmt.Println(output)
		return true
	}
	return false
}

func runCommand(command string, args ...string) (output string, status bool) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=amd64")

	var waitStatus syscall.WaitStatus
	combinedOutput, err := cmd.CombinedOutput()
	combinedOutputStr := string(combinedOutput)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			fmt.Printf("Output 1: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
			if waitStatus == 0 {
				return combinedOutputStr, true
			}
			return combinedOutputStr, false
		}
		return combinedOutputStr, false
	}
	// Success
	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	fmt.Printf("Output 2: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	if waitStatus == 0 {
		return combinedOutputStr, true
	}
	return combinedOutputStr, false
}
