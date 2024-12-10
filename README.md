# Objective
Creation of a URL Shortener

## How to Test the Code

Run the script from the terminal, specifying the port:

```bash
    ./shortner -p 2020
```
*Note: If No Port Is Specified, The Server Will Automatically Start On Port 8080.*




```bash
curl -X POST -d "url=http://www.google.it" http://localhost:2020/create
```
*Replace the url parameter with the link you want to convert into a short link.*


The server will respond with the new shortened URL.

Copy and paste the provided link into your browser.
