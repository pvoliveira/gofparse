![TRavis-ci](https://travis-ci.org/pvoliveira/gofparse.svg?branch=master) [![codecov](https://codecov.io/gh/pvoliveira/gofparse/branch/master/graph/badge.svg)](https://codecov.io/gh/pvoliveira/gofparse)

# gofparse
Simple file parser written in Golang

## Let's see that
The application just reads line-by-line a file and break into lines with fields, all these lines and fields must be passed by a file config in JSON like below:

```JSON
{
  "fileDescription":"test config",
  "options": [],
  "linesConfig": [
    {
      "description": "header",
      "identifierField": {
        "description": "type of record",
        "initPos": 0,
        "size": 2,
        "typeData": "string",
        "key": "00"
      },
      "fields": [
        {
          "description": "type of record",
          "initPos": 0,
          "size": 2,
          "typeData": "string"
        },
        {
          "description": "description of record",
          "initPos": 2,
          "size": 7,
          "typeData": "string"
        },
        {
          "description": "observation",
          "initPos": 10,
          "size": 20,
          "typeData": "string"
        }
      ]
    },
    {
      "description": "type 1",
      "identifierField": {
        "description": "type of record",
        "initPos": 0,
        "size": 2,
        "typeData": "string",
        "key": "01"
      },
      "fields": [
        {
          "description": "type of record",
          "initPos": 0,
          "size": 2,
          "typeData": "string"
        },
        {
          "description": "description of record",
          "initPos": 2,
          "size": 7,
          "typeData": "string"
        },
        {
          "description": "observation",
          "initPos": 10,
          "size": 20,
          "typeData": "string"
        }
      ]
    },
    {
      "description": "type phone",
      "identifierField": {
        "description": "type of record",
        "initPos": 0,
        "size": 2,
        "typeData": "string",
        "key": "02"
      },
      "fields": [
        {
          "description": "type of record",
          "initPos": 0,
          "size": 2,
          "typeData": "string"
        },
        {
          "description": "description of record",
          "initPos": 2,
          "size": 7,
          "typeData": "string"
        },
        {
          "description": "phone",
          "initPos": 10,
          "size": 11,
          "typeData": "string"
        },
        {
          "description": "observation",
          "initPos": 21,
          "size": 9,
          "typeData": "string"
        }
      ]

    }

  ]

}
```
