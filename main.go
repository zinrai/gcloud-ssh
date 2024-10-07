package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Defaults     Environment            `yaml:"defaults"`
	Environments map[string]Environment `yaml:"environments"`
}

type Environment struct {
	Project   string `yaml:"project"`
	Host      string `yaml:"host,omitempty"`
	Zone      string `yaml:"zone,omitempty"`
	User      string `yaml:"user,omitempty"`
	SocksPort int    `yaml:"socks_port,omitempty"`
}

func main() {
	envType := flag.String("env", "", "Environment type (required)")
	useSocks := flag.Bool("socks", false, "Use SOCKS proxy")
	configPath := flag.String("config", filepath.Join(os.Getenv("HOME"), ".config", "gcloud-ssh.yaml"), "Path to config file")

	flag.Parse()

	if *envType == "" {
		fmt.Println("Error: Environment type is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := checkGcloudCommand(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	env, ok := config.Environments[*envType]
	if !ok {
		fmt.Printf("Error: No match for environment type '%s'\n", *envType)
		os.Exit(1)
	}

	if env.Host == "" {
		env.Host = config.Defaults.Host
	}
	if env.Zone == "" {
		env.Zone = config.Defaults.Zone
	}
	if env.User == "" {
		env.User = config.Defaults.User
	}
	if env.SocksPort == 0 {
		env.SocksPort = config.Defaults.SocksPort
	}

	fmt.Printf("SSH Login to \"%s\" at %s\n\n", *envType, time.Now().Format(time.RFC1123))

	args := []string{
		"compute", "ssh",
		fmt.Sprintf("%s@%s-%s", env.User, env.Host, *envType),
		fmt.Sprintf("--project=%s", env.Project),
		fmt.Sprintf("--zone=%s", env.Zone),
	}

	if *useSocks {
		args = append(args, "--", "-N", "-p", "22", "-D", fmt.Sprintf("localhost:%d", env.SocksPort))
	}

	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println("Executing command:", cmd.String())

	err = cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("SSH connection terminated. Exit status: %d\n", exitError.ExitCode())
		} else {
			fmt.Printf("Error executing gcloud command: %v\n", err)
		}
	} else {
		fmt.Println("SSH connection closed normally.")
	}
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func checkGcloudCommand() error {
	_, err := exec.LookPath("gcloud")
	if err != nil {
		return fmt.Errorf("gcloud command not found in PATH. Please ensure it's installed and in your PATH")
	}
	return nil
}
