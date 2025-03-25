# Csv2SQLite

A simple command-line tool to convert CSV files into SQLite databases. This tool reads a CSV file and imports its
contents into a SQLite database table with appropriate column names.

## Features

- Read CSV files with header rows
- Automatically create SQLite tables with column names from CSV headers
- Insert all records from CSV to SQLite
- Simple command-line interface
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Pre-built Binaries

Pre-built binaries are available for Linux, macOS, and Windows from
the [releases page](https://github.com/khmlpy/csv2sqlite/releases).

### Build from Source

To build from source, you need Go 1.18 or later installed on your machine.

```shell
# Clone the repository
git clone https://github.com/khmlpy/csv2sqlite.git
cd csv2sqlite

# Build
cd src
go build -o csv2sqlite
```

## Usage

### Basic Usage

```shell
./csv2sqlite read --csv <path_to_csv> --table <table_name> --db <path_to_sqlite_db>
```

Or using short flags:

```shell
./csv2sqlite read -c <path_to_csv> -t <table_name> -d <path_to_sqlite_db>
```

### Example

```shell
./csv2sqlite read --csv sample.csv --table sample --db sample.sqlite
```

This will:

1. Read `sample.csv` file
2. Create a table named `sample` in the SQLite database `sample.sqlite`
3. Import all records from the CSV file into the database table

### Sample CSV Format

The CSV file should have a header row with column names. All data types are treated as TEXT in SQLite.

Example CSV:

```
"id","name","age"
"1","Alice","22"
"2","Bob","16"
"3","charlie","32"
```

## How It Works

1. The tool opens and reads the CSV file
2. It uses the first row of the CSV as column names
3. It creates a SQLite table with those column names (all columns are TEXT type)
4. It inserts each subsequent row from the CSV as a record in the table
5. It reports the number of records inserted upon completion

## Building for Different Platforms

Use the provided Makefile to build for different platforms:

```shell
# Build for Linux
make build-linux

# Build for macOS
make build-darwin

# Build for Windows
make build-windows

# Build all platforms and create zip archives
make build-release-artifacts
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
