# filebuilder - A quick uitility for creating file structures and comparing them

`go get github.com/jackmordaunt/filebuilder`

## Why?

I had trouble trying to mock out the filesystem for command line apps. 
Creating the files and comparing the files procedurally is tedious.
I wanted a declarative way to specify "these are the files I have", "these are the files I want", and "compare this directory with that directory".

To that end, I don't know if this is the 'best' way to do it, or even a good way to do it.

Code suggestions and project suggestions very welcome.

```go
func main() {

        // Declare the files you want, you can nest with relative paths or use 
        // a flat list, or both.
        entries := []filebuilder.Entry{
                filebuilder.File{Path: "foo/bar/baz.exe"},
                filebuilder.File{Path: "foo/baz.exe"},
                filebuilder.File{Path: "baz.exe"},
                filebuilder.Dir{Path: "bar", Entries: []filebuilder.Entry{
                        filebuilder.File{Path: "baz.txt", Content: []byte("foo")},
                        filebuilder.File{Path: "foo/bar/baz.txt", Content: []byte("bar")},
                }},
        }
        
        // Grab a filesystem implementation, or use the default by passing in nil.
        // The optional root will be the parent of the provided entries.
        fs := afero.NewMemMapFs()
        cleanup, err := filebuilder.Build(fs, "parent", entries...)
        if err != nil {
                log.Fatalf("failed creating entries: %v", err)
        }

        // Optional cleanup func which erases all files created. 
        defer func() {
                if err := cleanup(); err != nil {
                        log.Fatalf("failed cleanup of files: %v", err)
                }
        }()

        // fs is stateful, you can build up the file tree over multiple calls.
        _, err = filebuilder.Build(fs, "inside/this/folder", entries)
        if err != nil {
                log.Fatalf("failed creating entries: %v", err)
        }

        // Compare the directories.
        diff, ok, err := filebuilder.CompareDirectories(fs, "parent", "inside/this/folder")
        if err != nil {
                log.Fatalf("error while comparing directories: %v", err)
        }
        if !ok {
                log.Printf("directories are not equivalent, %v\n", diff)
        }
}
```