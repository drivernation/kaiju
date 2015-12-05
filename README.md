Kaiju
===
[![Build Status](https://travis-ci.org/drivernation/kaiju.svg?branch=master)](https://travis-ci.org/drivernation/kaiju)

Kaiju is a web application framework for the Go programming language. It brings together the standard Go libraries and a number of third party libraries in a sane way, providing developers with a single, easy-to-use web app framework.

Kaiju makes use of the following third party libraries:
* [Gorrilla Mux](https://github.com/gorilla/mux) for muxing HTTP(S) requests.
* [Seelog](https://github.com/cihub/seelog) for logging all the things.
* [Go-Yaml](https://github.com/go-yaml/yaml) for working with YAML files.
* And more to come!

Aside from the aforementioned third party libraries, Kaiju provides some of its own tooling, such as a service framework for convenient management of long running processes.

For more information on using Kaiju, checkout out our [wiki](https://github.com/drivernation/kaiju/wiki)!

Want to contribute to Dropwizard?
---
Before working on the code, if you plan to contribute changes, please read the following [CONTRIBUTING](CONTRIBUTING.md) document.

Need help or found an issue?
---
When reporting an issue through the [issue tracker](https://github.com/drivernation/kaiju/issues?state=open)
on GitHub, please use the following guidelines:

* Check existing issues to see if it has been addressed already
* The version of Kaiju you are using
* A short description of the issue you are experiencing and the expected outcome
* Description of how someone else can reproduce the problem
* Paste error output or logs in your issue or in a Gist. If pasting them in the GitHub
issue, wrap it in three backticks: ```  so that it renders nicely
* Write a unit test to show the issue!
