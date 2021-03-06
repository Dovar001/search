package search

import (
	"context"
	"os"
	"strings"
	"sync"
	"log"
)

//Result  описывает один результат поиска
type Result struct {
	//Фраза которую искали
	Phrase string
	//Целиком вся строка в котором нашли вхождение( без /n или /r/n в конце)
	Line string
	//Номер строки(начиная с 1) в котором нашли вхождение
	LineNum int64
	//Номер позиции(начиная с 1) в котором нашли вхождение
	ColNum int64
}

func FindAllMatchTextInFile(phrase, fileName string) (res []Result){
	read,_:=os.ReadFile(fileName)
	fstr:= string(read)
    filestr:=strings.Split(fstr,"\n")

	 if (len(filestr) > 0){

		filestr = filestr[:len(filestr)-1]
	 }
	 for i, line := range filestr {
        
		   if (strings.Contains(line,phrase)){
            
			result := Result{
				Phrase:  phrase,
				Line: line,
				LineNum: int64(i+1),
				ColNum: int64(strings.Index(line,phrase)) + 1, 

			}
			res = append(res, result) 
		}
	 }
return res
}

func FindAnyMatchTextInFile(phrase, filetext string) (res Result) {

	

	temp := strings.Split(filetext, "\n")

	for i, line := range temp {
		
		if strings.Contains(line, phrase) {

			return Result{
				Phrase:  phrase,
				Line:    line,
				LineNum: int64(i + 1),
				ColNum:  int64(strings.Index(line, phrase)) + 1,
			}

		}
	}

	return res
}


//All ищет все вхождение phrase в текстовых файлах files
func All(ctx context.Context, phrase string, files []string) <-chan []Result{

	ch:= make(chan []Result)
    
	wg:= sync.WaitGroup{}

	ctx,cancel:=context.WithCancel(ctx)

	for i := 0; i < len(files); i++ {
		wg.Add(1)

		go func(ctx context.Context, filename string, i int, ch chan<- []Result) {
           defer wg.Done()
		   res:=FindAllMatchTextInFile(phrase,filename)

		   if len(res) > 0 {
             ch <- res
		   }
			
		}(ctx, files[i], i, ch)
		
		go func() {
			defer close(ch)
			wg.Wait()
		  }()
		 
	}
	cancel()
 return ch	
}


func Any(ctx context.Context, phrase string, files []string) <-chan Result {

	ch := make(chan Result)
	wg := sync.WaitGroup{}
	result := Result{}
	for i := 0; i < len(files); i++ {

		data, err := os.ReadFile(files[i])
		if err != nil {
			log.Println("error not opened file err => ", err)
		}
		filetext := string(data)

		if strings.Contains(filetext, phrase) {
			res := FindAnyMatchTextInFile(phrase, filetext)
			if (Result{}) != res {
				result = res
				break
			}
		}

	}
	log.Println(result)

	wg.Add(1)
	go func(ctx context.Context, ch chan<- Result) {
		defer wg.Done()
		if (Result{}) != result {
			ch <- result
		}
	}(ctx, ch)

	go func() {
		defer close(ch)
		wg.Wait()
	}()
	return ch
}