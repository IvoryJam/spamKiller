package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Print("URL: ")
		buffer := bufio.NewReader(os.Stdin)
		url, err := buffer.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading URL from input: %s\n", err)
		}
		url = strings.TrimSpace(url)
		abuseEmail, err := getAbuseEmail(url)
		if err != nil {
			log.Fatalf("Error running whois: %s\n", err)
		}
		emailLink := emailBuilder(abuseEmail, url)
		fmt.Println(emailLink)
		exec.Command("flatpak", "run", "org.mozilla.firefox", emailLink).Run()

	} else {
		for _, i := range os.Args[1:] {
			abuseEmail, err := getAbuseEmail(i)
			if err != nil {
				log.Fatalf("Error running whois: %s\n", err)
			}
			emailLink := emailBuilder(abuseEmail, i)
			fmt.Println(emailLink)
			exec.Command("flatpak", "run", "org.mozilla.firefox", emailLink).Run()
		}
	}
}

func getAbuseEmail(url string) (abuseEmail string, err error) {
	output, err := exec.Command("whois", url).Output()
	outputString := string(output)
	emailRegex := regexp.MustCompile(`Registrar Abuse Contact Email: .*`)
	abuseEmail = emailRegex.FindString(outputString)
	if len(strings.Split(abuseEmail, ": ")) == 2 {
		abuseEmail = strings.Split(abuseEmail, ": ")[1]
	} else if len(strings.Split(abuseEmail, ": ")) == 1 {
		err = fmt.Errorf("Error getting abuse email: %s", url)
	}
	return
}

func emailBuilder(abuseEmail, spamURL string) (emailLink string) {
	emailLink = fmt.Sprintf("mailto:%s?subject=Abuse%%20from%%20%s&body=Hey%%20I%%20received%%20an%%20unsolicited%%20SMS%%20message%%20directing%%20me%%20to%%20%s", abuseEmail, spamURL, spamURL)
	return
}
