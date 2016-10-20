/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Deliver my data, Mr. Json!"
description = "Using JSON in RESTful Web API's in Go"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-10-20"
publishdate = "2016-10-20"
draft = "true"
domains = ["Internet And Web"]
tags = ["JSON", "REST", "AJAX"]
categories = ["Tutorial"]
+++

JSON is the *lingua franca* of exchanging data over the net and between applications written in different programming languages. In this article, we create tiny JSON server and client programs in Go.

<!--more-->

## What is JSON?

JSON is a standard data format for exchanging data over the net and between applications. It became popular as a more readable alternative to XML, with the added benefit that Javascript code can read and write JSON data out of the box. Still, almost any other popular programming language has at least one library for handling JSON data.

JSON is a human-readable text format, which makes it suitable for configuration files as well.

## So how does JSON data look like?

The syntax of JSON is almost like you would create a JavaScript object. This, by the way, is also where the name "JSON" comes from: It is an acronym for "Javascript Standard Object Notation".

In its base form, JSON data is a list of name-value entities. As an example, let's create some weather information in JSON.

```json
{
	"location": "Zzyzx",
	"weather": "sunny"
}
```

Here we have two entries, location and weather. Both are strings. Note that the names must be enclosed in double quotes. This is not a requirement for Javascript, only for JSON.

JSON knows some other data types besides strings: numbers, booleans, a null value, arrays, and objects. Let's add some of these to our weather data:
A numeric temperature, a boolean to tell whether the temperature is measured in celsius or fahrenheit, an array with the temperature forecast for the next three days, and an object that describes the wind in terms of direction and speed.

```json
{
	"location": "Zzyzx",
	"weather": "sunny",
	"temperature": 30,
	"celsius": true,
	"temp_forecast": [ 27, 25, 28 ],
	"wind": {
		"direction": "NW",
		"speed": 15
	}
}
```

## What is not in JSON?

Some data types or special values cannot be expressed in JSON. Most notably, there is no date type. Any date value is converted to and from an ISO-8601 date string. Also, Javascript's special values `NaN`, `Infinity`, and `-Infinity` are simply turned into a `null` value.

## JSON and Go

Go's standard library includes the package `encoding/json` that makes working with JSON a snap. With this package, we can map JSON objects to Go struct types, and convert data between the two. Let's examine this in code.

*/

//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Defining the structure of our weather data in Go is straightforward.
type weatherData struct {
	locationName string
	weather      string
	temperature  int
	celsius      bool
	tempForecast []int `json: temp_forecast`
	wind         windData
}

type windData struct {
	direction string
	speed     int
}

// For encoding the Go data as JSON, we use the Marshal function from `encoding/json`.
func encode(w weatherData) ([]byte, error) {
	weatherJson, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}
	return weatherJson, nil
}

// This wasn't too difficult, so let's try decoding JSON data.
func decode(weatherJson []byte) (weatherData, error) {
	// First, we need a weatherData struct to receive the decoded data.
	w := weatherData{}
	// Now we can pass the JSON data and a pointer to our struct to Unmashal.
	err := json.Unmarshal(weatherJson, &w)
	// Return the weatherData struct and any error that Unmarshal may have produced.
	return w, err

}

/*
Let's test this in a small server application.

The server's tasks is simple: It shall receive JSON-encoded location data and respond with weather information.

*/

// The location data is just a latitude and a longitude.
type loc struct {
	lat int
	lon int
}

//
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// We start by decoding the client request. The request consists of bare JSON location information,
	// so we can decode straight away into a location structure.
	loc, err := decode(r.Body)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	// Now it's time to prepare our response by setting up a weatherData structure.
	// We could try fetching the data from a weather service, but for the purpose of
	// demonstrating JSON handling, let's just use some mock-up data.
	weather := weatherData{
		locationName: "Zzyzx",
		weather:      "cloudy",
		temperature:  31,
		celsius:      true,
		tempForecast: []int{30, 32, 29},
		wind: windData{
			direction: "S",
			speed:     20,
		},
	}

	//

}

func main() {
	if len(os.Args) < 2 || (os.Args(1) != "server" && os.Args(1) != "client") {
		fmt.Println(`Usage:
json server
json client
`)
		os.Exit(1)
	}
	if os.Args(1) == "server" {
		if err := server(); err != nil {
			log.Fatal("Server error: ", err)
		}
	} else {
		if err := client(); err != nil {
			log.Fatal("Client error: ", err)
		}
	}
}

/*
There are a couple more things about JSON than what fits into the 5-10 minutes format of a screencast.

If you want to learn more about JSON, like decoding JSON data of an unknown structure, or how to implement stream encoding or decoding, the article "JSON and Go" from the official Go blog is a great place to start:

https://blog.golang.org/json-and-go

The transcript for this video is available here:

https://appliedgo.net/json


That's it for today, thanks for watching and happy coding!

*/
