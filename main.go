// Copyright 2020 John Cramb. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/jcramb/cedict"
)

func main() {
	var bSay = flag.Bool("say", false, "Enable Speaker Say")
	flag.Parse()
	homedir := os.Getenv("HOME")
	mdbname := path.Join(homedir, ".cache", "dict_mdb.txt.gz")
	d, err := cedict.Load(mdbname)
	if err != nil {
		d = cedict.New()
	}
	s := strings.Join(flag.Args(), " ")
	var src string
	if len(s) > 0 {
		if cedict.IsHanzi(s) {
			fmt.Printf("[input] hanzi\n")

			// convert to pinyin
			src = cedict.PinyinTones(d.HanziToPinyin(s))
			fmt.Printf("%s\n", src)
			if *bSay {
				spdsay(src)
			}
		} else {
			fmt.Printf("[input] english \n")

			// search by meaning
			elements := d.GetByMeaning(s)
			for _, e := range elements {
				line := e.Marshal()
				fmt.Printf("%s\n", line)

				re := regexp.MustCompile(`\[([^\[\]]*)\]`)
				submatchall := re.FindAllString(line, -1)
				for _, element := range submatchall {
					element = strings.Trim(element, "[")
					element = strings.Trim(element, "]")
					if *bSay {
						spdsay(element)
					}
				}
			}
		}
	}

	d.Save(mdbname)
}

func spdsay(str string) {
	str = strings.ReplaceAll(str, ":", "")
	cmd := exec.Command("/usr/bin/espeak-ng", "-g10", "-a100", "-s160", "-v", "cmn-latn-pinyin", str)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %v\n", err)
	}
}
