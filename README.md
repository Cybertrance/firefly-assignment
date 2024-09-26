# **WebWordRank**

## 📖 **Description**

**WebWordRank** is a concurrent Go-based application designed to process and analyze web content, such as articles and essays. It utilizes efficient concurrency patterns, rate limiting, and a modular architecture to fetch and process content from multiple sources. The application is built with scalability and performance in mind, leveraging Go's powerful goroutines for fast and efficient processing.

The project includes CI/CD pipelines for automated testing, building, and deployment of cross-platform binaries for Linux and Windows.

## ✨ **Features**

- **Concurrent Processing**: Efficiently fetches and processes web content in parallel using Go's goroutines.
- **Rate Limiting**: Includes configurable, built-in rate limiting to avoid overwhelming external services with too many requests.
- **Cross-Platform Support**: Builds binaries for both Linux and Windows.
- **CI/CD Integration**: Automated testing, building, and deployment pipelines using GitHub Actions.
- **Customizable**: Includes configuration options to configure aspects of the application.

## 🚀 **Installation**

### Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/dl/) 1.21 or higher
- [Git](https://git-scm.com/)

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/Cybertrance/firefly-assignment.git
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Build the application (On windows):

   ```bash
   go build -o firefly.exe
   ```

4. Run the application (On windows):
   ```bash
   ./firefly.exe
   ```

Alternatively, you can also download Linux or Windows binaries from [the latest CI pipeline.](https://github.com/Cybertrance/firefly-assignment/actions/runs/11054683629)

## 🧪 **Running Tests**

This project includes a comprehensive test suite. To run all the unit tests:

```bash
go test ./... -v
```

Tests are automatically executed in the CI pipeline on every commit and pull request to ensure code quality and functionality.

## 📦 **Cross-Platform Builds**

You can manually build the application for different platforms using the following commands:

- **Linux**:

  ```bash
  GOOS=linux GOARCH=amd64 go build -o [your-repo-name] .
  ```

- **Windows**:
  ```bash
  GOOS=windows GOARCH=amd64 go build -o [your-repo-name].exe .
  ```

Alternatively, binaries for both platforms are automatically built using GitHub Actions and are available as build artifacts in the CI/CD pipeline.

## 🚀 **CI/CD Pipeline**

This project uses **GitHub Actions** to automate the testing, building, and deployment processes:

- **Testing**: All commits and pull requests trigger the test suite to ensure that the code passes all unit tests before proceeding.
- **Build**: If the tests pass, the application is built for both Linux and Windows platforms. The binaries are then uploaded as artifacts.
- **Static File Management**: Static files are automatically copied to the build directory and included in the build artifacts.
- **Cross-Platform Support**: The CI pipeline generates platform-specific binaries, ensuring compatibility with both Linux and Windows environments.

The CI/CD pipeline is defined in the `.github/workflows` directory and includes steps for:

- Running tests
- Building Linux and Windows binaries
- Uploading build artifacts

## 🔧 **Configuration**

You can configure the application through the `config.yaml` file.

### Example `config.yaml` file:

| **Configuration**         | **Default Value**                                                         | **Description**                                            |
| ------------------------- | ------------------------------------------------------------------------- | ---------------------------------------------------------- |
| `top_results`             | `10`                                                                      | Number of top results to display after processing content. |
| `source_url_filename`     | `"endg-urls"`                                                             | Filename that contains the list of URLs for scraping.      |
| `word_bank_url`           | `"https://raw.githubusercontent.com/dwyl/english-words/master/words.txt"` | # URL to fetch a word bank.                                |
| `container_selector`      | `".caas-body"`                                                            | CSS selector used to target the content in HTML scraping.  |
| `requests_per_second`     | `20`                                                                      | Maximum number of requests allowed per second.             |
| `burst_size`              | `20`                                                                      | Maximum burst size allowed when rate limiting requests.    |
| `max_concurrent_requests` | `20`                                                                      | Maximum number of requests that can be made concurrently.  |
| `max_retries`             | `3`                                                                       | Number of retries allowed when requests fail.              |
| `maxRedirects`            | `5`                                                                       | Maximum number of redirects that are followed per request. |

## 📜 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
