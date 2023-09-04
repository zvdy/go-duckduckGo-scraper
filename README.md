# DuckDuckGo Web Scraper

This is a simple command-line web scraper written in Go that scrapes search results from DuckDuckGo.

## Features

- Search for a specified term on DuckDuckGo.
- Customize the number of search results to retrieve.
- Save the search results to a text file.
- Automatically generate a User-Agent header for HTTP requests.

## Usage

### Prerequisites

- Go installed on your system.

### Installation

1. Clone this repository to your local machine:

   ```sh
   git clone https://github.com/yourusername/duckduckgo-web-scraper.git
    ```
    

1. Change to the project directory:

    
    ```sh
    cd duckduckgo-web-scraper
    ```
1. Build the Go program:

    ```sh
    go build
    ```
1. Run the program:

    ```sh
    ./duckduckgo-web-scraper -search "your search term" -count 10
    ```
    >Replace "your search term" with the term you want to search for and -count with the number of search results you want to retrieve (default is 10).

### Flags

- search or -s: Specify the search term.

- user-agent: Customize the User-Agent string for HTTP requests (optional).

- output: Specify the output file name (optional). By default, it's named after the search term.

- count: Specify the number of search results to retrieve (default is 10).

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


### Acknowledgments

- __github.com/PuerkitoBio/goquery:__ Used for parsing HTML content.

- __DuckDuckGo:__ The search engine whwere the search results are retrieved from.