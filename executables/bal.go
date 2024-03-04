package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	scriptPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println("Error getting script path:", err)
		os.Exit(1)
	}
	version := os.Getenv("VERSION")
	//fmt.Print(scriptPath)
	//scriptPath = "/usr/lib/ballerina/bin/bal"
	var javaCommand string
	javaCommand = "java"
	scriptDir := filepath.Dir(scriptPath)

	switch osType := runtime.GOOS; osType {
	case "darwin":
		// path is already set
	default:
		scriptDir, err = filepath.EvalSymlinks(scriptDir)
		if err != nil {
			fmt.Println("Error resolving symlinks:", err)
			os.Exit(1)
		}
		scriptDir, err = filepath.Abs(scriptDir)
		if err != nil {
			fmt.Println("Error getting absolute path:", err)
			os.Exit(1)
		}
	}

	// Check for specific Java versions
	javaVersions := []string{"jdk-17.0.7+7-jre", "jdk-11.0.18+10-jre", "jdk-11.0.15+10-jre", "jdk-11.0.8+10-jre", "jdk8u265-b01-jre"}
	for _, version := range javaVersions {
		javaDir := filepath.Join(scriptDir, "..", "dependencies", version)
		if _, err := os.Stat(javaDir); err == nil {
			javaCommand = filepath.Join(javaDir, "bin", "java")
			break
		}
	}
	//bal completion bash or bal completion zsh commands implementation
	if len(os.Args) >= 3 && os.Args[1] == "completion" {
		completionScriptPath := filepath.Join(scriptDir, "..", "scripts", "bal_completion.bash")

		if _, err := os.Stat(completionScriptPath); err == nil {
			if os.Args[2] == "zsh" {
				fmt.Printf("autoload -U +X bashcompinit && bashcompinit\n")
				fmt.Printf("autoload -U +X compinit && compinit\n\n")
				fmt.Printf("#!/usr/bin/env bash\n\n")
				printCompletionScript(completionScriptPath)
			} else if os.Args[2] == "bash" {
				fmt.Printf("#!/usr/bin/env bash\n\n")
				printCompletionScript(completionScriptPath)
			} else {
				fmt.Printf("ballerina: unknown command '%s'\n", os.Args[2])
				os.Exit(1)
			}
		} else {
			fmt.Println("Completion scripts not found")
			os.Exit(1)
		}
		os.Exit(0)
	}
	jarFileName := "ballerina-command-" + version + ".jar"
	javaCommandPath := filepath.Join(scriptDir, "..", "lib", jarFileName)
	runCommand := false
	runBallerina := true
	if len(os.Args) > 1 {
		if os.Args[1] == "dist" || os.Args[1] == "update" || (os.Args[1] == "dist" && os.Args[2] == "update") {
			runCommand = true
			runBallerina = false
		}
		if os.Args[1] == "build" {
			runCommand = true
		}
		if runCommand {
			if os.Args[1] == "build" {
				//javaCommandPath := filepath.Join(scriptDir, "..", "lib", "ballerina-command-1.4.2.jar") //ballerina-command-@version@.jar
				args := append([]string{"-jar", javaCommandPath}, "build")
				executeCommand(javaCommand, args)
			} else {
				if runtime.GOOS == "darwin" {
					os.Setenv("BALLERINA_MAC_ARCHITECTURE", runtime.GOARCH)
				}
				//javaCommandPath := filepath.Join(scriptDir, "..", "lib", "ballerina-command-1.4.2.jar")
				args := append([]string{"-jar", javaCommandPath}, os.Args[1:]...)
				executeCommand(javaCommand, args)
			}

			if os.Args[1] == "update" {
				tmpDir := filepath.Join(scriptDir, "..", "ballerina-command-tmp")
				if _, err := os.Stat(tmpDir); err == nil {
					cmd := exec.Command(filepath.Join(tmpDir, "install"))
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						fmt.Println("Update failed due to errors")
						_ = os.RemoveAll(tmpDir)
						if exitErr, ok := err.(*exec.ExitError); ok {
							os.Exit(exitErr.ExitCode())
						} else {
							//fmt.Println("Error running jar file: ", err)
							os.Exit(1)
						}

					}
					_ = os.RemoveAll(tmpDir)
					fmt.Println("Update successfully completed")
					fmt.Println()
					fmt.Println("If you want to update the Ballerina distribution, use 'bal dist update'")
					os.Exit(0)

				}
			}
		}
	}

	if runBallerina {
		//determining version of ballerina
		distributionFilePath := filepath.Join(scriptDir, "..", "distributions", "ballerina-version")
		ballerinaVersion, err := readVersionFromFile(distributionFilePath)
		if err != nil {
			//fmt.Println("Error reading distribution version:", err)
			os.Exit(1)
		}

		// Check and read user-specific version file
		userVersionFilePath := filepath.Join(os.Getenv("HOME"), ".ballerina", "ballerina-version")
		//fmt.Println(userVersionFilePath)
		userBallerinaVersion, err := readVersionFromFile(userVersionFilePath)
		if err != nil {
			//fmt.Println("Error reading user-specific version:", err)
			os.Exit(1)
		}
		// Override with user-specific version if exists
		if userBallerinaVersion != "" && checkDirectoryExists(filepath.Join(scriptDir, "..", "distributions", userBallerinaVersion)) {
			ballerinaVersion = userBallerinaVersion

		}
		ballerinaVersion = strings.TrimSuffix(ballerinaVersion, "\n")
		//fmt.Println("Selected Ballerina Version:", ballerinaVersion)
		//Setting  Ballerina home and excectuion of bal file
		ballerinaHome := filepath.Join(scriptDir, "..", "distributions", ballerinaVersion)
		os.Setenv("BALLERINA_HOME", ballerinaHome)
		balPath := filepath.Join(ballerinaHome, "bin", "ball")
		ballerinaPath := filepath.Join(ballerinaHome, "bin", "ballerina")
		if _, err := os.Stat(balPath); err == nil {
			executeCommand(balPath, os.Args[1:])
		} else {
			// Check if 'ballerina' executable exists
			if _, err := os.Stat(ballerinaPath); err == nil {
				executeCommand(ballerinaPath, os.Args[1:])
			} else {
				fmt.Println("Distribution does not exist, use 'bal dist pull <version>'")
				os.Exit(1)

			}
		}
		if len(os.Args) > 1 {
			if os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" ||
				os.Args[1] == "version" || os.Args[1] == "-v" || os.Args[1] == "--version" || (os.Args[1] == "help" && os.Args[2] == "") {
				//jarFilePath := filepath.Join(scriptDir, "..", "lib", "ballerina-command-1.4.2.jar")
				args := append([]string{"-jar", javaCommandPath}, os.Args[1:]...)
				executeCommand(javaCommand, args)
			}
		} else if len(os.Args) == 1 {
			//jarFilePath := filepath.Join(scriptDir, "..", "lib", "ballerina-command-1.4.2.jar")
			args := append([]string{"-jar"}, javaCommandPath)
			executeCommand(javaCommand, args)
		} else {
			exitCode := 1
			os.Exit(exitCode)
		}
	}
	os.Exit(0)
}

func printCompletionScript(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Failed to generate the completion script", err)
		os.Exit(1)
	}
	fmt.Print(string(content))
}
func readVersionFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func checkDirectoryExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	return err == nil && info.IsDir()
}

func executeCommand(commandPath string, args []string) {
	cmd := exec.Command(commandPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}
