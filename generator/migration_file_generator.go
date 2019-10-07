package generator

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	//Migration define migration db file generator
	Migration interface {
		Generate(name string) error
	}

	migrate struct{}
)

//NewMigrationGenerator :nodoc:
func NewMigrationGenerator() Migration {
	return &migrate{}
}

//GenerateMigration to generate sql migration file
func (m migrate) Generate(name string) error {
	err := m.checkMigrationFolderExists()
	if err != nil {
		fmt.Println(err)
	}
	migrationFile := []byte(`-- +migrate Up notransaction` + "\n\n" + `-- +migrate Down`)
	migrationFileName := "db/migration/" + m.createUniqueTime() + "_" + strings.ToLower(name) + ".sql"
	err = ioutil.WriteFile(migrationFileName, migrationFile, 0666)
	if err != nil {
		return err
	}
	fmt.Println(migrationFileName + " created")
	return nil
}

func (m migrate) checkMigrationFolderExists() error {
	_, err := os.Stat("db/migration/")
	if os.IsNotExist(err) { //check if db/migration folder is already exist
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("db/migration/ folder not found, want to  create (Y/N)? ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return errors.New("fail when read user input")
		}
		ans := strings.Contains(input, "y")
		if ans {
			err = m.createMigrationFolder()
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func (m migrate) createMigrationFolder() error {
	_, err := os.Stat("db/")
	if os.IsNotExist(err) { //check if db folder is already exist
		err := os.Mkdir("db/", 0666)
		if err != nil {
			return err
		}
	}
	err = os.Mkdir("db/migration/", 0666)
	if err != nil {
		return err
	}
	return nil
}

func (m migrate) createUniqueTime() string {
	now := time.Now()
	splitDate := strings.Split(now.Format("01/02/2006"), "/") // mm/dd/yyyy
	newDate := splitDate[2] + splitDate[0] + splitDate[1]
	hr, min, sc := now.Clock()
	hour := strconv.Itoa(hr)
	minute := strconv.Itoa(min)
	sec := strconv.Itoa(sc)
	if len(hour) == 1 {
		hour = "0" + hour
	}
	if len(minute) == 1 {
		minute = "0" + minute
	}
	if len(sec) == 1 {
		sec = "0" + sec
	}

	return newDate + hour + minute + sec
}