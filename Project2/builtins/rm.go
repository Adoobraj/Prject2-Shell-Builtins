//moo0065
package builtins

import (
   "fmt"
   "os"
)

func RemoveFile(args ...string) error {
   
   if len(args) == 0 {
      return fmt.Errorf("usage: rm <filename>")
   }
   
   filename := args[0]
   err := os.Remove(filename)
   if err != nil{
      return err
   }
   
   return nil
}

 
