package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/ypapax/logrus_conf"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	if err := func() error {
		log.SetFlags(log.LstdFlags | log.Llongfile)
		if err := logrus_conf.PrepareFromEnv("linesep"); err != nil {
			return errors.WithStack(err)
		}
		if len(os.Args) < 3 {
			execName := "linesep"
			if len(os.Args) > 0 {
				execName = os.Args[0]
			}
			return errors.Errorf("usage: %+v filePath lineNumber", execName)
		}
		file := os.Args[1]
		logrus.Infof("reading file %+v   ...", file)
		b, err := os.ReadFile(file)
		if err != nil {
			return errors.WithStack(err)
		}
		lineNumberStr := os.Args[2]
		lineNumberFromOne, err := strconv.Atoi(lineNumberStr)
		if err != nil {
			return errors.WithStack(err)
		}
		lineNumberFromZero := lineNumberFromOne - 1 // count from 1
		lines := strings.Split(string(b), "\n")
		if len(lines) <= lineNumberFromZero {
			return errors.Errorf("not enough lines %+v for getting line %+v", len(lines), lineNumberFromOne)
		}
		l := lines[lineNumberFromZero]
		lineFileNameToSave := "/tmp/" + path.Base(file) + "_" + lineNumberStr + "_" + path.Ext(file)
		logrus.Infof("saving line '%+v' to %+v", l, lineFileNameToSave)
		if err := os.WriteFile(lineFileNameToSave, []byte(l), 0666); err != nil {
			return errors.WithStack(err)
		}
		logrus.Infof("saved line '%+v' to %+v", l, lineFileNameToSave)

		return nil
	}(); err != nil {
		log.Printf("%+v", err)
		logrus.Errorf("error: %+v", err)
	}

}
