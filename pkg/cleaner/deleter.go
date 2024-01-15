package cleaner

import (
	"bufio"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
)

func Delete(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		log.Err(err, "Can't open file.")
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log.Info(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Err(err, "Scanner error")
	}

	return err
}
