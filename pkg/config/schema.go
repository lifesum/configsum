package config

import "github.com/xeipuuv/gojsonschema"

const schemaDefBaseCreate = `
{
  "$schema":"http://json-schema.org/draft-06/schema#",
  "title":"Base update",
  "description":"Request data for base config updates.",
  "type":"object",
  "required":[
    "client_id", "name"
  ],
  "properties": {
    "client_id": {
      "type": "string"
    },
    "name": {
      "type": "string"
    }
  }
}`

const schemaDefBaseUpdate = `
{
  "$schema":"http://json-schema.org/draft-06/schema#",
  "title":"Base update",
  "description":"Request data for base config updates.",
  "type":"object",
  "required":[
    "parameters"
  ],
  "properties":{
    "parameters":{
      "type":"object",
      "additionalProperties":false,
      "minProperties":1,
      "patternProperties":{
        "^feature_([a-z-]+)_([a-z]+)$":{
          "anyOf":[
            {
              "type":"boolean"
            },
            {
              "type":"number"
            },
            {
              "type":"string"
            }
          ]
        }
      }
    }
  }
}`

const schemaDefUserRender = `
{
  "$schema": "http://json-schema.org/draft-06/schema#",
  "title": "Render context",
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
      "type": [ "null", "object" ],
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
          "type": "object",
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

var (
	schemaBaseCreateRequest *gojsonschema.Schema
	schemaBaseUpdateRequest *gojsonschema.Schema
	schemaUserRenderRequest *gojsonschema.Schema
)

func init() {
	var err error

	schemaBaseCreateRequest, err = gojsonschema.NewSchema(
		gojsonschema.NewStringLoader(schemaDefBaseCreate),
	)
	if err != nil {
	}

	schemaBaseUpdateRequest, err = gojsonschema.NewSchema(
		gojsonschema.NewStringLoader(schemaDefBaseUpdate),
	)
	if err != nil {
		panic(err)
	}

	schemaUserRenderRequest, err = gojsonschema.NewSchema(
		gojsonschema.NewStringLoader(schemaDefUserRender),
	)
	if err != nil {
		panic(err)
	}
}
