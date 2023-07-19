package main

import (
	"bufio"
	"fmt"
	///"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

// writeLines writes the lines to the given file.

func runClone(wg *sync.WaitGroup , url string)  {

	wg.Add(1)
	defer wg.Done()

	os.Chdir("./results")
	
	cmd := exec.Command("git" , "clone" , url)

    _ , err := cmd.Output()

    if err != nil {
        log.Fatal(err)
    }

	fmt.Printf("%v Cloned Successfully !\n" , url)
}

func runPull(wg *sync.WaitGroup , subDir string)  {

	wg.Add(1)
	defer wg.Done()

	os.Chdir("./results" + "/" + subDir)
	
	cmd := exec.Command("npm" , "i")

    res , err := cmd.Output()

    if err != nil {
        log.Fatal(err)
    }

	fmt.Printf("\033[36m%v \033[0m Pulled Successfully !\n \033[0m" , subDir)
	fmt.Printf("\033[0m Msg: \033[33m%v \n \033[0m" , string(res))
}

func main() {

	var wg sync.WaitGroup;

	urlsForscan , err := readLines("../crawler/results.csv");
	if err != nil {
		log.Fatal(err)
	}

	var hasResults bool;

	if _, err := os.Stat("./results"); err != nil {
		if os.IsNotExist(err) {
			hasResults = false
			os.Mkdir("results" ,0755)
		}
	}else{
		hasResults = true;
	}

	dirs , err := os.ReadDir("./results")
	if err != nil {
		log.Fatal(err)
	}

	if(hasResults && len(dirs) != 0){
		for _ , dir := range dirs {
			dirName := dir.Name()
				go runPull(&wg , dirName)
		}

	}else{
		// for _ , url := range urlsForscan{
		// 	var innerUrl = url
		// 	go func(){
		// 		runClone(&wg , innerUrl)
		// 	}()
		// }	
		go runClone(&wg , urlsForscan[0])
		go runClone(&wg , urlsForscan[1])
		go runClone(&wg , urlsForscan[2])
		go runClone(&wg , urlsForscan[4])

	}

	wg.Wait()
}