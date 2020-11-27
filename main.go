package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

// Non-capitalized names to avoid exporting struct.
type contact struct {
	name  string
	email string
	phone string
}

var contactsList []contact

func main() {
	readContactsFile("./contacts.txt")
	var response string
	for {
		fmt.Print("> ")
		response = strings.ToLower(readSentence())
		processInput(response)
	}
}

// Read the contacts.txt file to import contacts, then return slice of contacts.
func readContactsFile(filename string) {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	// Will crash if line is longer than 64kb, more than enough for this program
	scanner := bufio.NewScanner(file)
	// Creates a slice instead of an array with no specified length!
	// Pretty cool feature in Go, although there's a lot of debate
	// about which you should use. Slices are built on top of arrays
	// and they're generally more common than arrays in Go code.
	// for is used as a while loop in Golang
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " | ")
		contactInfo := contact{line[0], line[1], line[2]}
		contactsList = append(contactsList, contactInfo)
	}
	// Make sure everything went well
	check(scanner.Err())
}

func processInput(input string) {
	switch input {
	case "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "add contact":
		addContact()
	case "show all contacts":
		for _, contact := range contactsList {
			fmt.Printf("%v's contact info:\n", contact.name)
			fmt.Printf("email: %v\n", contact.email)
			fmt.Printf("phone: %v\n\n", contact.phone)
		}
	default:
		if strings.HasPrefix(input, "show contacts with") {
			stringFields := strings.Fields(input)
			printContact(stringFields[3], strings.Join(stringFields[4:], " "))
		} else {
			fmt.Println("Huh?")
		}
	}
}

func addContact() {
	fmt.Print("name: ")
	name := readSentence()
	fmt.Print("email: ")
	email := readSentence()
	fmt.Print("phone: ")
	phone := readSentence()
	contactsList = append(contactsList, contact{name, email, phone})
	saveContacts("./contacts.txt")
}

func printContact(field, key string) {
	for _, contact := range contactsList {
		if strings.ToLower(getField(&contact, field)) == key {
			fmt.Printf("%v's contact info:\n", contact.name)
			fmt.Printf("email: %v\n", contact.email)
			fmt.Printf("phone: %v\n", contact.phone)
		}
	}
}

func saveContacts(filename string) {
	file, err := os.Create(filename)
	check(err)
	defer file.Close()
	var line string
	for _, contact := range contactsList {
		line = contact.name + " | " + contact.email + " | " + contact.phone + "\n"
		_, err := file.WriteString(line)
		check(err)
	}
}

// HELPER FUNCTIONS

// https://stackoverflow.com/questions/18930910/access-struct-property-by-name
func getField(v *contact, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

func readSentence() string {
	var strInput string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		strInput = scanner.Text()
	}
	return strInput
}

// https://gobyexample.com/reading-files
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
