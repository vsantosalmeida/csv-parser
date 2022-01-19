# CSV Parser

 Is a Golang library to process csv files and normalize the data to an employee entity. 

## How to run the project

First execute the build of the bin file.
```bash
make build
```

Copy yours csv files to the bin directory.
```bash
cp *.csv ./bin
```

Run the binary with the -files params (each file separated by `,`).
```bash
cd bin
./csv-parser.bin -f=roster1.csv,roster2.csv
```

For each file will require an input with the name of the columns.
Case your file don't have the column, just hit enter.

IMAGE

After the execution, if the files are processed with success one or both of that files will be created with the results.

**employee-{timestamp}.json**

**badData-{timestamp}.json**

Check coverage (will open in your browser the code coverage.)
```bash
make test cover-html
```

## Architecture

For this project, I choose to use a simple way to process the files, which is receiving from the input the structure of the CSV files. 
Thus ensuring in a better way that the parser will be able to interpret each file in the best way. 
As architecture, I use a pattern that in my opinion creates a better pattern in the structure of the project and ensures a more scalable code for any new feature. 
The clean architecture, which is a pattern to develop software independent of frameworks, UI, or any external technologies.


[Clean Architecture by Elton Minetto](https://eltonminetto.dev/en/post/2018-03-05-clean-architecture-using-go/)

## Tradeoffs

- Processing multiple files ends up being a lot of work
- As the name of the columns needs to be input, it makes possible many typos.

## What I would evolve?

- A better way to translate the names of the columns might create a custom flag implementation to parse an arg to a map.
- Process each file asynchronously with multiple results files.
- A flexible result file output, may be received by parameter the desirable file like CSV, JSON, or XML, and if not passed use JSON as default.
- A flexible parser, that can process CSV, JSON, or XML files.
