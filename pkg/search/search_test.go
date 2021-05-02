package search

import (
	"context"
	"log"
	//"go/types"
	"testing"
)

func BenchmarkAll(b *testing.B) {

  ctx:= context.Background()
 files := [] string{
	"../../data/file.txt", 
 }
want := [] Result{
	{
	Phrase:"Auto",
	Line:"40a6238e-ccf0-42d9-adf4-364c540c2d30;1;2;Auto;INPROGRESS",
	LineNum:2,
	ColNum:42,
	},
}
	for i := 0; i < b.N; i++ {
		ch:=All(ctx, "Auto",files)
		result:= <- ch

		if result[i] == want[i] {
			b.Fatalf("invalid result, got %v", result)
		}
	}
}



func TestAll(t *testing.T){


	ctx:= context.Background()
	files := [] string{
	   "../../data/file.txt", 
	}

	ch:=All(ctx, "Auto",files)

	s,ok := <-ch

	if !ok {
		t.Errorf("fuction All error +> %v", ok)
	}

	log.Println("---------", s)
}