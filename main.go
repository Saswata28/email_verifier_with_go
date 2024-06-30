package main

import(
	"fmt"
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"regexp"
)

func isValidEmail(email string) bool {
    const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9\-]+\.[a-zA-Z]{2,}(\.[a-zA-Z]{2,})*$`
    re := regexp.MustCompile(emailRegex)
    return re.MatchString(email)
}

func extractDomain(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return ""
	}
	return email[atIndex+1:]
}

func checkDomain(domain string){
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)
	if err != nil{
		log.Printf("\033[31mError: %v\033[0m\n", err)
	}
	if len(mxRecords) > 0{
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil{
		log.Printf("\033[31mError: %v\033[0m\n", err)
	}
	
	for _, record := range txtRecords{
		if strings.HasPrefix(record, "v=spf1"){
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil{
		log.Printf("\033[31mError: %v\033[0m\n", err)
	}

	for _, record := range dmarcRecords{
		if strings.HasPrefix(record, "v=DMARC1"){
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("Your domain: %v\n hasMX: %v\n hasSPF: %v\n spfRecord is: %v\n hasDMARC: %v\n dmarcRecord is: %v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)

	if (hasMX == true && hasSPF == false){
		fmt.Printf("\033[33mThis is a valid email domain, but it's not secure because it doesn't have SPF, so it's not recommended to use this domain for mail services.\033[0m\n")
	}
	if (hasMX == true && hasSPF == true && hasDMARC == true){
		fmt.Printf("\033[32mThis is a valid and secure domain for mail services\033[0m\n")
	}
	if (hasMX == false){
		fmt.Printf("\033[31mThis is not a valid email domain.\033[0m\n")
	}
}

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Email checker, enter email (type 'exit' or 'ctrl+c' to quit): ")
	for scanner.Scan(){
		input := strings.TrimSpace(scanner.Text())
		if input == "exit" {
			fmt.Printf("Exiting...\n")
			break
		}
		if isValidEmail(input){
			domain := extractDomain(input)
			checkDomain(domain)
			fmt.Printf("Email checker, enter email (type 'exit' to quit): ")
		} else {
			fmt.Printf("\033[31mInvalid Email\033[0m\n")
		}
		
	}

	if err := scanner.Err(); err != nil{
		log.Fatal("\033[31mError: could not read from input: %v\n", err)
	}
	
}
// \033[31m    \033[0m