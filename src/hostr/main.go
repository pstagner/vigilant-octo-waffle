package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// HostEntry represents an entry in the hosts file.
type HostEntry struct {
	IP   string
	Domain string
}

// loadEnv loads environment variables from a .env file.
func loadEnv(filePath string) (map[string]string, error) {
	env := make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening .env file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		env[key] = value
	}
	return env, nil
}

// getOperatingSystem returns the operating system name.
func getOperatingSystem() string {
	return runtime.GOOS
}
// getHostsFilePath returns the path to the hosts file based on the operating system.
func getHostsFilePath(osName string) string {
	switch osName {
	case "windows":
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
	return "/etc/hosts"
	case "darwin":
		return "/etc/hosts"
	default:
		log.Fatal("Unsupported operating system")
		return ""
}
}

// readEnvFile reads the .env file and returns a map of environment variables.
func readEnvFile(filePath string) (map[string]string, error) {
	env := make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening .env file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		env[key] = value
	}
	return env, nil
}

// processTemplate replaces variables in a template string with values from the environment.
func processTemplate(template string, env map[string]string) string {
	re := regexp.MustCompile(`\${(\w+)}`)
	return re.ReplaceAllStringFunc(template, func(match string) string {
		key := match[2 : len(match)-1]
		if value, ok := env[key]; ok {
			return value
		}
		return match // Return the original match if the variable is not found
	})
}

// addEntryToHostsFile adds an entry to the hosts file.
func addEntryToHostsFile(hostsFilePath string, ip string, domain string) error {
	hostsContent, err := ioutil.ReadFile(hostsFilePath)
	if err != nil {
		return fmt.Errorf("error reading hosts file: %w", err)
	}
	entry := fmt.Sprintf("%s\t%s\n", ip, domain)
	if strings.Contains(string(hostsContent), entry) {
		return nil // Entry already exists
	}
	hostsContent = append(hostsContent, []byte(entry)...)
	return ioutil.WriteFile(hostsFilePath, hostsContent, 0644)
}

// readHostsFile reads the contents of the hosts file.
func readHostsFile(hostsFilePath string) (string, error) {
	return ioutil.ReadFile(hostsFilePath)
}

// writeHostsFile writes the contents to the hosts file.
func writeHostsFile(hostsFilePath string, content string) error {
	return ioutil.WriteFile(hostsFilePath, []byte(content), 0644)
}

func main() {
	osName := getOperatingSystem()
	hostsFilePath := getHostsFilePath(osName)

	env, err := loadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}

	ip := env["IP"]
	domain := env["DOMAIN"]

	if ip == "" || domain == "" {
		log.Fatal("IP and DOMAIN must be set in .env file")
	}

	err = addEntryToHostsFile(hostsFilePath, ip, domain)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Entry added to hosts file successfully!")
}






















	entry := fmt.Sprintf("%s\t%s\n", ip, domain)

	if strings.Contains(string(hostsContent), entry) {
		return nil // Entry already exists
	}


	hostsContent = append(hostsContent, []byte(entry)...)
	return ioutil.WriteFile(hostsFilePath, hostsContent, 0644)
}

func main() {
	env, err := loadEnv(".env")
	if err != nil {


		log.Fatal(err)
	}
	ip := env["IP"]
	domain := env["DOMAIN"]
	if ip == "" || domain == "" {


		log.Fatal("IP and DOMAIN must be set in .env file")
	}


	err = addEntryToHostsFile(hostsFilePath, ip, domain)
	if err != nil {


		log.Fatal(err)
	}

	fmt.Println("Entry added to hosts file successfully!")
}

