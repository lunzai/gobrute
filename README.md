# Gobrute - RESTful API Brute-Force Login Script

## Description
**Gobrute** is a RESTful API brute-forcing tool created to practice ethical hacking, specifically for testing login passwords. The script is simple yet effective for learning, such as with the [OWASP Juice Shop](https://owasp.org/www-project-juice-shop/) admin password cracking challenge. Created as an alternative to Burp Suite Community (due to rate limitations) and Hydra (due to documentation complexity), Gobrute provides basic parameter customization, multi-threading, and clear output, making it easy for users with intermediate programming skills to expand or adjust for other use cases.

> **Note**: This tool is strictly for educational purposes. Use it only in safe, ethical, and legal environments with permission.

## Introduction
Gobrute uses Go’s concurrency and customizable request features to brute-force API login endpoints. Users can specify a wordlist, username, failure response string, timeout, thread count, and request payload format, allowing quick and effective login attempts for RESTful applications.

Some parameters are inspired by the Hydra project, and we extend our thanks to Hydra’s developers for the initial parameter ideas.

### Features
- **Threaded Execution**: Multi-threading with adjustable concurrency speeds up brute-forcing.
- **Customizable Payloads**: Define request payloads easily, adapting to various RESTful login endpoints.
- **Progress Tracking**: A simple progress bar shows the number of attempts made.
- **Error Handling**: Gobrute handles intermittent errors smoothly, allowing the script to continue running.

## Usage
### Prerequisites
- Install [Go](https://go.dev/dl/) to build and run Gobrute.
- Set up a test environment such as [OWASP Juice Shop](https://owasp.org/www-project-juice-shop/) for safe hacking practice.

### Installation
Clone this repository and navigate into the directory:
```bash
git clone https://github.com/lunzai/gobrute.git
cd gobrute
```

Initialize the Go module and install dependencies:
```bash
go mod init gobrute
go get -u github.com/schollz/progressbar/v3
```

Build the program:
```bash
go build -o gobrute
```

### Parameters and Usage
Run the program with the required parameters:

```bash
./gobrute -u <URL> -w <wordlist> -l <username> -f <failure_message> -p <payload_format> -t <threads> -o <timeout>
```

#### Parameter Details
- `-u` **(Required)**: Target URL for the API login endpoint.
- `-w` **(Required)**: Path to the password wordlist file.
- `-l` **(Required)**: Login username to test.
- `-f` **(Required)**: String indicating a failed login response. If this string appears in the response, the attempt is marked as unsuccessful.
- `-p` **(Required)**: Payload format, where `^USER^` will be replaced with the login username and `^PASS^` with each password attempt.
- `-t` **(Optional)**: Number of concurrent threads (default: 10).
- `-o` **(Optional)**: Request timeout in milliseconds (default: 2000 ms).

### Example Command
```bash
./gobrute -u http://juice.test/rest/user/login -w passwords.txt -l admin@juice-sh.op -f "Invalid email or password." -p '{"email":"^USER^","password":"^PASS^"}' -t 10 -o 2000
```

This command runs the brute-force attempt using:
- `passwords.txt` as the wordlist,
- `admin@juice-sh.op` as the username,
- with 10 concurrent threads,
- and a 2000ms timeout for each attempt.

### Progress and Results
- The script displays a progress bar indicating the number of attempts made.
- Upon finding a valid login, the script immediately outputs the successful username and password combination and stops further attempts.
- If no successful login is found after all attempts, it displays a failure message.

## Limitations
1. **Limited Error-Handling Scope**: While Gobrute continues on most connection errors, it lacks advanced recovery from complex server-side issues or non-standard error handling.
2. **Rate-Limiting**: This script doesn’t account for rate-limiting imposed by the server, which may cause temporary blocking on certain services.
3. **Testing Only for JSON Payloads**: The script is built specifically for JSON payloads, though users can modify it to support other formats.
4. **Single Username at a Time**: The script accepts only one username per run; for multiple usernames, run the script separately for each user.

## Acknowledgements
Special thanks to the creators of [Hydra](https://github.com/vanhauser-thc/thc-hydra) for the initial inspiration behind the parameters used in this script.

**Disclaimer**: This tool is intended for educational purposes only. Unauthorized use against any system is illegal and unethical. Ensure you have permission before testing any system.
