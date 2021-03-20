# MonGolang
## `MongoDB`
## ` +  Golang`
## `            = Mon Golang`
### Goal

Simplify Go programs which use MongoDB.
Create a library which supports the 10% of the calls which are needed 90% of the time, while easily allowing the use of the other 90% of MongoDB capabilities which are needed 10% of the time.

1. Support the 10% of the calls which are used 90% of the time
2. The other 90% of the calls which are used 10% of the time should be no more difficult to do then without this library
3. Calls should resemble the calls made via the MongoDB Console
4. Although error checking is important, it shouldn't get in the way of being able to chain calls such as:
   
    ```go
      db.Coll("someCollection").
        Find(`{"lastName":"Johnson"}`).
        Sort(`{"firstName":1}`).
        Limit(5).
        Pretty()
    ```

5. Don't extend the capabilities of the MongoDB Golang driver, just simplify the use of what's there
6. New MongoDB releases should not require changes to the library or programs that use the library.

### With Jupyter Notebook
Combined with  [Jupyter Notebook and GopherNotes](https://github.com/gopherdata/gophernotes), this provides a simple interface to run MongoDB Shell like commands using Go. Not only easy to iterate and easy to rerun, but also a great way to prototype code which can then be copied and pasted into a Go program. Below is an example query in Juptyter notebook using JSON strings.

![Example 01 Using MonGolang in Jupyter Notebook](misc/MonGolang_V02.1._Test01.png?raw=true)
![Example 02 Using MonGolang in Jupyter Notebook](misc/MonGolang_V02.1._Test02.png?raw=true)


[See here](https://docs.google.com/presentation/d/1zq8-n0w0uiy9AIK9kaOiZgIL6VEmUc1FBDpbImZ4RLw/edit?usp=sharing) for a brief slide presentation on running MonGolang on Jupyter and then using the code in a standalone compiled Go program.

There is a PDF demo of the use of MonGolang in Jupyter in this repo at [MonGolang_v0.2.1_Demo.pdf](https://github.com/ddgarrett/mongolang/blob/main/MonGolang_v0.2.1_Demo.pdf) 



