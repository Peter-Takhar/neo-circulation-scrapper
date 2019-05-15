# NEO Total Circulation Scraper

This repository contains a simple way to start a NEO node and to obtain the total circulating supply of NEO.

First, make sure you have Docker installed on your Debian System.

To run the NEO node as a daemon, run:

`docker run -it -d -p 10332:10332 petertakhar/ubuntu-neo-cli`

To run the scrapper program, run:

`go run neo.go`
