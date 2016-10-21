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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// Defining the structure of our weather data in Go is straightforward.
// Note that the json package only encodes public fields of a struct.
// The JSON fields are all lowercase, so we need to map the struct field
// names to the corresponding JSON field names.
// Luckily, [Go structs come with a string tag feature](https://golang.org/ref/spec#Struct_types).
// This way we can tag every struct field with the corresponding JSON field name.
type weatherData struct {
	LocationName string   `json: locationName`
	Weather      string   `json: weather`
	Temperature  int      `json: temperature`
	Celsius      bool     `json: celsius`
	TempForecast []int    `json: temp_forecast`
	Wind         windData `json: wind`
}

type windData struct {
	Direction string `json: direction`
	Speed     int    `json: speed`
}

// Let's implement a tiny server application. The client sends its location, and
// the server responds by sending weather data.
//
// Location data is just a latitude and a longitude.
type loc struct {
	Lat float32 `json: lat`
	Lon float32 `json: lat`
}

// For the server, we need a function for handling the request
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// First, we need a location struct to receive the decoded data.
	location := loc{}

	// The location data is inside the request body which is an io.ReadCloser,
	// but we need a byte slice for unmarshalling.
	// ReadAll from ioutil just comes in handy.
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}
	// Now we can decode the request data using the Unmarshal function.
	err = json.Unmarshal(jsn, &location)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	// To see if the request was correctly received, let's print it to the console.
	log.Printf("Received: %v\n", location)

	// Now it's time to prepare our response by setting up a weatherData structure.
	// We could try fetching the data from a weather service, but for the purpose of
	// demonstrating JSON handling, let's just use some mock-up data.
	weather := weatherData{
		LocationName: "Zzyzx",
		Weather:      "cloudy",
		Temperature:  31,
		Celsius:      true,
		TempForecast: []int{30, 32, 29},
		Wind: windData{
			Direction: "S",
			Speed:     20,
		},
	}

	// For encoding the Go struct as JSON, we use the Marshal function from `encoding/json`.
	weatherJson, err := json.Marshal(weather)
	spew.Dump(weatherJson)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}
	// We send a JSON response, so we need to set the Content-Type header accordingly.
	w.Header().Set("Content-Type", "application/json")

	// Sending the response is as easy as writing to the ResponseWriter object.
	w.Write(weatherJson)

}

// Thanks to Go's http package, starting the server is a piece of cake.
func server() {
	http.HandleFunc("/", weatherHandler)
	http.ListenAndServe(":8080", nil)
}

// Our mock client is even simpler than the server.
func client() {
	// Again we create JSON by marshalling a struct; in this case a loc struct literal.
	locJson, err := json.Marshal(loc{Lat: 35.14326, Lon: -116.104})
	spew.Dump(locJson)
	// Then we set up a new HTTP request for posting the JSON data to local port 8080.
	req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer(locJson))
	req.Header.Set("Content-Type", "application/json")

	// An HTTP client will send our HTTP request to the server and collect the response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Finally, we print the received response and close the response body.
	// Note we use ReadAll here for simplicity. Be careful when using ReadAll in larger
	// projects, as reading large files can consume a lot of memory.
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Response: ", string(body))
	resp.Body.Close()
}

// Function main calls either the client() function or the server() function,
// depending on the first parameter on the command line.
func main() {
	if len(os.Args) < 2 || (os.Args[1] != "server" && os.Args[1] != "client") {
		fmt.Println(`Usage:
json server
json client
`)
		os.Exit(1)
	}
	if os.Args[1] == "server" {
		server()
	} else {
		client()
	}
}

/*
There are a couple more things about JSON than what fits into the 5-10 minutes format of a screencast.

If you want to learn more about JSON, like decoding JSON data of an unknown structure, or how to implement stream encoding or decoding, the article "JSON and Go" from the official Go blog is a great place to start:

https://blog.golang.org/json-and-go

The transcript for this video is available at

https://appliedgo.net/json

where you also can find the link to the source code as well as installation instructions.

That's it for today, thanks for watching and happy coding!

- - -

## Getting and installing the code

As always, use -d to prevent the binary from showing up in your path.

First, get the code:

    go get -d github.com/appliedgo/json.go

Then, compile the code and start the server:

    cd $GOPATH/src/github.com/appliedgo/json
	go build
	./json server

Finally, open another shell and run the client:

    cd $GOPATH/src/github.com/appliedgo/json
	./json client

You should then see the server response in the second shell.
*/
