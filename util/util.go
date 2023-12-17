package util

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-password/password"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

func IsCommandExist(cmd string) bool {
	_, err := exec.LookPath(cmd)

	return err == nil
}

func Sed(old, new, filePath string) error {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	m1 := regexp.MustCompile(old)

	fileString := string(fileData)
	fileString = m1.ReplaceAllString(string(fileData), new)
	fileData = []byte(fileString)

	err = ioutil.WriteFile(filePath, fileData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func AppendFile(filePath string, data string) {
	f, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data); err != nil {
		log.Println(err)
	}
}

func FileGetContents(filename string) string {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		os.Exit(1)
	}

	return string(fileData)
}

func GeneratePassword(length int) string {
	res, err := password.Generate(length, 10, 0, false, false)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func LoadEnvFile(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}

func MakeSSHKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	// generate and write private key as PEM
	var privKeyBuf strings.Builder

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&privKeyBuf, privateKeyPEM); err != nil {
		return "", "", err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	var pubKeyBuf strings.Builder
	pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))

	return pubKeyBuf.String(), privKeyBuf.String(), nil
}

func GetCurrentDir() string {
	if currentDir, err := os.Getwd(); err == nil {
		return currentDir
	}

	return ""
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ExecCommand(cmdName string) {
	cmdArgs := strings.Fields(cmdName)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

func CreateFileIfNotExists(filename string) {
	if !FileExists(filename) {
		os.Create(filename)
	}
}

func Mkdir(dirs []string, perm os.FileMode) {
	for i := 0; i < len(dirs); i++ {
		os.Mkdir(dirs[i], perm)
	}
}

func Chmod(name string, mode os.FileMode) {
	err := os.Chmod(name, mode)
	if err != nil {
		log.Fatal(err)
	}
}

func DownloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Contains(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}
