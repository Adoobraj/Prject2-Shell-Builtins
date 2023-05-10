package builtins

import (
   "fmt"
   "os"
)

func main() {
   // Delete the file named "file.txt".
   err := os.Remove("file.txt")
   if err != nil {
       fmt.Println(err)
       return
   }
   fmt.Println("File successfully deleted")
}

 
