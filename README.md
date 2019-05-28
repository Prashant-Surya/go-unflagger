### Go-UnFlagger

The main purpose of this project is to remove feature flags in code wherever present by just running a single command.

Currently this supports removal of flags by two conditions:

##### 1. Removing Flag by name:

While running the command if the flag's name is provided, it checks entire file to identify where a flag by this name is present and remove the block

##### 2. Removing Flag's that are past the feature launch date:

This case each feature flag is written with a separator(__) in the name and when flagger type is chosen as date it'll parse through all flags and check if it's date is less than current date and remove those.

Currently existence of flag is determined by checking if there is an `if` statement with condition similar to `<some_package>.FeatureFlags.<flag_name>`. We can make this flexible by taking the condition regex as command-line argument. 

For this project to work with your code base use a struct with feature flags as the attributes.

### Usage:

Download the latest release from [https://github.com/Prashant-Surya/go-unflagger/releases](https://github.com/Prashant-Surya/go-unflagger/releases) and place it in a folder available in `$PATH`

```
-date-format string
	Format of the date embedded in flag (default "2006_01_02" -> "yyyy_mm_dd")
-name string
	Name of the flag to be removed
-path string
	Relative or Absolute Path of the file or directory
-recursive
	Recursively parse flags. (Enable in case -path is a directory)
-type string
	Flagger Type. Possible values date, name (default "date")
-write
	Enable this flag to update contents to file or it'll be written to stdout
```	

##### Examples:

1. Removing a feature flag with name `launch_v1` from a directory at `/Users/surya/go/src/testProject`.
```
flagger -path /Users/surya/go/src/testProject -recursive -type name -name launch_v1
```

2.  Removing a feature flag with name `launch_v1` from a file at `/Users/surya/go/src/testProject/main.go`.
```
flagger -path /Users/surya/go/src/testProject/main.go -type name -name launch_v1
```

3. Removing multiple feature flags with name `launch_v1__2019_04_04`, `launch_v2__2019_05_04` from a file at `/Users/surya/go/src/testProject/main.go` using date type flagger.
```
flagger -path /Users/surya/go/src/testProject/main.go -type date
```


### How it works:

When a .go file is parsed using `ast` library it returns a tree with root node containing all the statements in that file as it's children.

For example the following code generates the following ast.

#### Code: 
```
package main

import (
	"config"
	"fmt"
)	

func main() {
	if config.FeatureFlags.Enable {
		fmt.Println("Enabled")
	} else {
		fmt.Println("Disabled")
	}
}
```

#### AST: 

![pre](https://github.com/Prashant-Surya/go-unflagger/blob/master/resources/ast1.png)

After running:

```
flagger -path /Users/surya/go/src/testProject/main.go -type name -name Enable
```

AST node of the if statement is replaced with the body of the `IF` statment. This works even when there are multiple statements under an if.


#### Code:
```
package main

import (
	"config"
	"fmt"
)

func main() {
	fmt.Println("Enabled")
}
```

#### AST:

![post](https://github.com/Prashant-Surya/go-unflagger/blob/master/resources/ast2.png)
