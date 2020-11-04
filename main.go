/*
 * OS3LABGEN - Dead simple pdf to dokuwiki lab template generator for OS3 students
 * Copyright (c) 2020 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"code.sajari.com/docconv"
)

const version = "v1.0"

func main() {

	// check args
	if len(os.Args) < 2 {
		fmt.Println("Dead simple tool to convert a PDF document with OS3 lab exercises into a dokuwiki template.")
		fmt.Println("Provide the pdf as first argument, the resulting dokuwiki markup will be written to stdout.")
		fmt.Println()
		fmt.Println("    usage: os3labgen <lab.pdf>")
		fmt.Println()
		fmt.Println("error: please provide a pdf as first argument")
		os.Exit(1)
	}

	// extract text from pdf
    res, err := docconv.ConvertPath(os.Args[1])
    if err != nil {
        log.Fatal(err)
	}

	// print header
	fmt.Println("==== Lab Template for", os.Args[1], "====")
	fmt.Println()
	fmt.Println("> generated at", time.Now().UTC())
	fmt.Println("> with https://github.com/dreadl0ck/os3labgen", version)
	fmt.Println()
	
	// setup a few variables for state keeping
	var (
		taskStarted = false
		nlCount = 0
		foundAbstract = false
		previousTaskComplete = false
	)

	// for debugging
	//fmt.Println(ansi.Red, res.Body, ansi.Reset)

	// iterate over pdf text line by line
    for _, line := range strings.Split(res.Body, "\n") {
		
		// stop when hitting refs
		if strings.HasPrefix(line, "References") {
			// done
			break
		}

		// include abstract text after the header
		if strings.HasPrefix(line, "Abstract") {
			foundAbstract = true
			continue
		}
		if foundAbstract {
			if strings.HasPrefix(line, "Task") {
				foundAbstract = false
				fmt.Print("\n=== ", line)
				taskStarted = true
				continue
			}
			fmt.Println(line)
			continue
		}

		// for debugging
		//fmt.Println(ansi.Red, line, ansi.Reset)
		//fmt.Println(ansi.Yellow, taskStarted, nlCount, ansi.Reset)

		// process tasks
		if line == "" {

			// count newlines
			nlCount++
			
			// two consequtive newlines end a task description
			if nlCount == 1 {
				taskStarted = false
				nlCount = 0
				continue
			}
		}

		// task has begun
		if taskStarted {
			// check if line belongs to a task or starts a new one
			if !strings.HasPrefix(line, "Task") {
				if len(line) > 0 {
					fmt.Print(" ", line)
					previousTaskComplete = false
				}
			} else {
				taskStarted = false
				nlCount = 0
				fmt.Println(" ===", "\n")
				previousTaskComplete = true
			}
		}

		// check if line is the start of a task
		if strings.HasPrefix(line, "Task") && len(line) > 0 {

			// check if the previous task line was closed properly
			if !previousTaskComplete {
				fmt.Println(" ===", "\n")
				previousTaskComplete = true
			}

			// start new task
			fmt.Print("\n=== ", line)
			taskStarted = true
		}		
	}

	// close the last task
	fmt.Println(" ===")
}