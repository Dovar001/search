package search

import (
	"context"
	"os"
	"strings"
	"sync"
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

func FindAnyMatchTextInFile(phrase, fileName string) (res Result){
	read,_:=os.ReadFile(fileName)
	fstr:= string(read)
    filestr:=strings.Split(fstr,"\n")

	 if (len(filestr) > 0){

		filestr = filestr[:len(filestr)-1]
	 }
	 for i, line := range filestr {
        
		   if (strings.Contains(line,phrase)){
            
			res = Result{
				Phrase:  phrase,
				Line: line,
				LineNum: int64(i+1),
				ColNum: int64(strings.Index(line,phrase)) + 1, 

			}
		}
		break
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


func Any(ctx context.Context, phrase string, files []string) <-chan Result{

	ch:= make(chan Result)
    
	wg:= sync.WaitGroup{}

	ctx,cancel:=context.WithCancel(ctx)

	for i := 0; i < len(files); i++ {
		wg.Add(1)
      test:=Result{}
		go func(ctx context.Context, filename string, i int, ch chan<- Result) {
           defer wg.Done()
		   res:=FindAnyMatchTextInFile(phrase,filename)
 
          if res !=test{
			   ch <- res
			  cancel()    
		  }
 	     
		 }(ctx, files[i], i, ch)
		}	
		 go func() {
			defer close(ch)	
			wg.Wait()
			cancel()
		  }()
		  
	
	return ch
}