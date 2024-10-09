# go-yaml-diff

[![Go Version](https://img.shields.io/badge/go-1.23%2B-blue)](https://golang.org/doc/go1.23)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

`go-yaml-diff` is a command-line tool written in Go that allows you to compare two YAML files and highlight the semantic differences between them. This tool is designed to help developers and system administrators identify discrepancies in YAML configurations quickly and easily.

## Features
- **Semantic Comparison:** Ignores the order of fields while comparing YAML files, ensuring that logically identical files with different orderings do not appear as different.
- **Colored Output:** Highlights differences in a visually appealing way using color coding.
- **Easy-to-Read Context:** Displays lines of context around the differences for better understanding.
- **Simple CLI Usage:** Easy-to-run commands with a comprehensive `Makefile` for common tasks.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Makefile Commands](#makefile-commands)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Installation
Make sure you have Go installed on your system (version 1.18 or higher). If not, download it from the official [Go website](https://golang.org/dl/).

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jefersonwurz/go-yaml-diff.git
   cd go-yaml-diff
