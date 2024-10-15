# Golang Reporting Engine

## Intoduction

Golang Reporting Engine is a simple tool built using Go that allows you to generate reports based on Excel templates. You can pass JSON formated data, and the tool will populate the data into the specified cells of the Excel file, generating a downloadable report. Thanks to [Excelize](https://github.com/qax-os/excelize) for make it this project easier to implement.

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

Example JSON Data of "data":

```
{
    "Sheet1": {
        "A1": {
            	"id" : "12"
		"value": "this",
        },
        "A2": {
            	"id" : "12"
		"value": "is",
        },
        "A3": {
            	"id" : "12"
		"value": "hello",
        },
        "A4": {
            	"id" : "12"
		"value": "world",
        }
    },
    "Sheet2": {
        "A1": {
            	"id" : "12"
		"value": "data",
        },
        "A2": {
            	"id" : "12"
		"value": "for",
        },
        "A3": {
            	"id" : "12"
		"value": "report",
        },
        "A4": {
            	"id" : "12"
		"value": "engine",
        }
    }
}
```

In this example, the JSON defines two sheets (Sheet1 and Sheet2). Each sheet has several cells (A1, A2, etc.) with corresponding values to be written into the Excel file.

Example JSON Data of "style" to apply styling to cell with id = "12":

```
{
    "12": {
        "font": {
            "bold": true,
            "italic": false,
            "size": 16,
            "color": "#FFFFFF"
        },
        "fill": {
            "type": "pattern",
            "color": ["#4CAF50"],
            "pattern": 1
        },
        "border": [
            {
                "type": "left",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "right",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "top",
                "color": "#000000",
                "style": 1
            },
            {
                "type": "bottom",
                "color": "#000000",
                "style": 1
            }
        ],
        "alignment": {
            "horizontal": "center",
            "vertical": "center"
        }
    }
}
```
