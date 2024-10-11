# Golang Reporting Engine

## Intoduction

Golang Reporting Engine is a simple tool built using Go that allows you to generate reports based on Excel templates. You can pass JSON formated data, and the tool will populate the data into the specified cells of the Excel file, generating a downloadable report. Thanks to [Excelize](https://github.com/golang/go/issues/61881) for make it this project easier to implement.

# Main Functionality

This tool's main purpose is to take an existing Excel file template, populate it with provided JSON data, and allow the user to download the filled-out file. The process involves:

Reading an Excel template.
Populating the specified cells with data from the passed JSON file.
Saving the result as a new Excel file.

# Passed Parameters

You can pass the following parameters in the POST request:

name : The name of the output file to be downloaded.
file : The template file name (e.g., Book1.xlsx) used as the base for the report.
data : A JSON file containing the data to be populated in the Excel file. The JSON structure should follow a specific format (see example below).

# JSON Data Format Example

The JSON data format is expected to follow a structure of nested objects where each sheet in the Excel file is represented as an object, and each cell is a key-value pair. Each value object should contain the following:

value: The actual value to be inserted in the Excel cell.
id: An optional ID for the value (used for identification purposes).

Example JSON Data:

```
{
    "Sheet1": {
        "A1": {
            "value": "this",
            "id": 12
        },
        "A2": {
            "value": "is",
            "id": 12
        },
        "A3": {
            "value": "hello",
            "id": 12
        },
        "A4": {
            "value": "world",
            "id": 12
        }
    },
    "Sheet2": {
        "A1": {
            "value": "data",
            "id": 34
        },
        "A2": {
            "value": "for",
            "id": 34
        },
        "A3": {
            "value": "report",
            "id": 34
        },
        "A4": {
            "value": "engine",
            "id": 34
        }
    }
}
```

In this example, the JSON defines two sheets (Sheet1 and Sheet2). Each sheet has several cells (A1, A2, etc.) with corresponding values to be written into the Excel file.
