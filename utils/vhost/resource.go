// Copyright 2017 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vhost

import (
	"bytes"
	"io/ioutil"
	"net/http"

	frpLog "github.com/charlesbao/frpc/utils/log"
	"github.com/charlesbao/frpc/utils/version"
)

var (
	NotFoundPagePath = ""
)

const (
	NotFound = `<!DOCTYPE html>
<html>
<head>
	<title>Not Found</title>
	<style type="text/css">
	html {
		font-size: 62.5%;
	}
	body {
		background-color: #fff;
		color: #000;
		font-family: helvetica, arial, sans-serif;
		font-size: 1.4em;
		line-height: 1.5;
	}
	.centered {
		position: fixed;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
	}
	.emoji {
		font-size: 9em;
		text-align: center;
	}
	.title {
		font-size: 3em;
		text-align: center;
		line-height: 0em;
		color: grey;
	}
	.text {
		text-align: center;
	}
	</style>
</head>
<body class="centered">
	<div class="emoji">😭</div>
	<p class="title">Ooooops!</p>
	<p class="text">404 NOT FOUND.</p>
</body>
</html>
`
)

func getNotFoundPageContent() []byte {
	var (
		buf []byte
		err error
	)
	if NotFoundPagePath != "" {
		buf, err = ioutil.ReadFile(NotFoundPagePath)
		if err != nil {
			frpLog.Warn("read custom 404 page error: %v", err)
			buf = []byte(NotFound)
		}
	} else {
		buf = []byte(NotFound)
	}
	return buf
}

func notFoundResponse() *http.Response {
	header := make(http.Header)
	header.Set("server", "frp/"+version.Full())
	header.Set("Content-Type", "text/html")

	res := &http.Response{
		Status:     "Not Found",
		StatusCode: 404,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     header,
		Body:       ioutil.NopCloser(bytes.NewReader(getNotFoundPageContent())),
	}
	return res
}

func noAuthResponse() *http.Response {
	header := make(map[string][]string)
	header["WWW-Authenticate"] = []string{`Basic realm="Restricted"`}
	res := &http.Response{
		Status:     "401 Not authorized",
		StatusCode: 401,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     header,
	}
	return res
}
