# ASCII Art Web

## Description

ASCII Art Web is a web application that converts text into ASCII art using different banner styles. Built with Go, it provides a graphical user interface for generating ASCII art from user input. The application supports three different banner styles: standard, shadow, and thinkertoy.

Users can input any text (including multiple lines), select their preferred banner style, and receive beautifully formatted ASCII art output displayed directly on the webpage.

## Authors

- clintoo

## Usage: How to Run

### Prerequisites

- Go 1.16 or higher installed on your system

### Running the Application

1. Clone or navigate to the project directory:

```bash
cd /workspaces/test
```

2. Run the server:

```bash
go run .
```

3. Open your web browser and navigate to:

```
http://localhost:8080
```

4. Use the web interface:
   - Select an ASCII art style (standard, shadow, or thinkertoy)
   - Enter the text you want to convert
   - Click "Generate" to create your ASCII art
   - The result will be displayed below the form

5. To stop the server, press `Ctrl+C` in the terminal

## Implementation Details: Algorithm

### Banner File Format

Each banner file (`standard.txt`, `shadow.txt`, `thinkertoy.txt`) contains ASCII art representations for all printable ASCII characters (32-126):

- Line 1: Empty line (skipped during parsing)
- Lines 2-856: Character definitions (95 characters × 9 lines each)
- Each character consists of 8 lines of ASCII art followed by 1 separator line

### Core Algorithm

1. **Banner Loading** (`loadBanner` function):
   - Read the banner file
   - Skip the first empty line
   - Parse character definitions into a map `map[rune][]string`
   - Each character (space through tilde) maps to 8 lines of ASCII art
   - Formula: Character at ASCII code `n` starts at line `(n - 32) × 9`

2. **Text Rendering** (`render` function):
   - Normalize line endings (handle both `\r\n` and `\n`)
   - Split input text into lines
   - For each line:
     - If empty: output a single newline
     - If non-empty: render 8 rows of ASCII art
       - For each row (0-7):
         - For each character in the input line:
           - Append the corresponding row from the character's ASCII art
         - Add newline after each row
   - Return the complete ASCII art output

3. **HTTP Handlers**:
   - `GET /`: Serves the main page with the input form
   - `POST /ascii-art`: Processes form data and returns ASCII art

4. **Error Handling**:
   - 200 OK: Successful request
   - 400 Bad Request: Invalid input (non-printable characters or invalid banner)
   - 404 Not Found: Invalid route or missing template/banner file
   - 500 Internal Server Error: Unexpected server errors

### Project Structure

```
.
├── main.go                 # Entry point
├── go.mod                  # Go module file
├── README.md               # This file
├── ascii-art/
│   ├── asciiart.go        # ASCII art generation logic
│   └── banners/
│       ├── standard.txt   # Standard banner font
│       ├── shadow.txt     # Shadow banner font
│       └── thinkertoy.txt # Thinkertoy banner font
├── handler/
│   └── handler.go         # HTTP request handlers
├── server/
│   └── server.go          # Server setup and routing
├── static/
│   └── style.css          # CSS styling
└── templates/
    ├── layout.html        # Base layout template
    ├── index.html         # Main page template
    ├── notFound.html      # 404 error template
    └── serverError.html   # 500 error template
```

### Key Features

- Clean separation of concerns (MVC-like architecture)
- Robust error handling with appropriate HTTP status codes
- Input validation (only printable ASCII characters)
- Responsive web design with dark theme
- Three distinct ASCII art styles
- Support for multi-line text input
