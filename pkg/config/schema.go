package config

import "github.com/xeipuuv/gojsonschema"

const requestCapabilities = `
{
  "$schema": "http://json-schema.org/draft-06/schema#",
  "title": "Client payload",
  "description": "Common set of information to determine device capabilities and user provided info.",
  "type": "object",
  "properties": {
    "app": {
      "type": "object",
      "properties": {
        "version": {
          "description": "The version of the client application.",
          "type": "string"
        }
      },
      "required": [
        "version"
      ]
    },
    "metadata": {
      "type": "object",
      "additionalProperties": {
        "anyOf": [
          {
            "type": "string"
          },
          {
            "type": "integer"
          },
          {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          {
            "type": "array",
            "items": {
              "type": "integer"
            }
          }
        ]
      }
    },
    "device": {
      "type": "object",
      "properties": {
        "location": {
          "type": [ "null", "object" ],
          "properties": {
            "locale": {
              "description": "The device's locale setting according to ISO 639-3 and/or BCP 47.",
              "type": "string"
            },
            "timezoneOffset": {
              "description": "time offset from GMT",
              "type": "integer"
            }
          },
          "required": [
            "locale",
            "timezoneOffset"
          ]
        },
        "os": {
          "type": "object",
          "properties": {
            "platform": {
              "description": "The client platform that makes this request.",
              "enum": [
                "Android",
                "iOS",
                "WatchOS"
              ]
            },
            "version": {
              "description": "Version of the os that runs on the client platform",
              "type": "string"
            }
          },
          "required": [
            "platform",
            "version"
          ]
        }
      },
      "required": [
        "location",
        "os"
      ]
    },
    "user": {
      "type": "object",
      "properties": {
        "age": {
          "description": "Age of the application's logged in user.",
          "type": "integer"
        }
      }
    }
  },
  "required": [
    "app",
    "device"
  ]
}`

var decodeClientPayloadSchema *gojsonschema.Schema

func init() {
	var err error

	decodeClientPayloadSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(requestCapabilities))
	if err != nil {
		panic(err)
	}
}
