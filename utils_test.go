package em

import (
	"fmt"
	"testing"
)

func TestMessageWithLineNum_Advanced(t *testing.T)  {
	call()
}

func call()  {
	fmt.Println(MessageWithLineNumAdvanced("Test",1,10))
}