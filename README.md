# Objective
Creation of a URL Shortener

## How to Test the Code

1. Run the script from the terminal, specifying the port:
   ```bash
   ./shortner -p 2020
Note: If no port is specified, the server will automatically start on port 8080.
```
2. Execute a curl command like the following from the terminal:

 ```bash
curl -X POST -d "url=http://www.google.it" http://localhost:2020/create
```
Replace the url parameter with the link you want to convert into a short link.


4. The server will respond with the new shortened URL.

5. Copy and paste the provided link into your browser.
