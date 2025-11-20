# 🚀 Online Go Interpreter

<div align="center">

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![WebAssembly](https://img.shields.io/badge/WebAssembly-654FF0?style=for-the-badge&logo=webassembly&logoColor=white)
![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)
![License](https://img.shields.io/badge/License-ISC-blue?style=for-the-badge)

**A blazing-fast, browser-based Go code interpreter powered by WebAssembly**

[Live Demo](https://joshua-beatty.github.io/online-go-interpreter) • [Report Bug](https://github.com/Joshua-Beatty/online-go-interpreter/issues) • [Request Feature](https://github.com/Joshua-Beatty/online-go-interpreter/issues)

</div>

---

## ✨ Features

- **🌐 100% Client-Side** — Runs entirely in your browser with zero network calls after initial load
- **⚡ WebAssembly Powered** — Leverages WASM for near-native performance
- **🎨 Monaco Editor** — Industry-standard code editor with syntax highlighting and IntelliSense
- **🔒 Secure & Private** — Your code never leaves your browser
- **📦 No Installation Required** — Just open and start coding
- **🎯 Powered by Yaegi** — Uses the [Yaegi Go interpreter](https://github.com/traefik/yaegi) for reliable Go execution

---

## 🎯 Quick Start

### Prerequisites

Before you begin, ensure you have the following installed:
- [Node.js](https://nodejs.org/) (v14 or higher)
- [Go](https://golang.org/) (v1.16 or higher)
- npm (comes with Node.js)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Joshua-Beatty/online-go-interpreter.git
   cd online-go-interpreter
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Build and run**
   ```bash
   npm run dev
   ```

4. **Open your browser**
   
   Navigate to `http://localhost:8080` (or the port shown in your terminal)

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────┐
│           Browser (Client-Side)             │
├─────────────────────────────────────────────┤
│                                             │
│  ┌───────────────────────────────────────┐ │
│  │      Monaco Editor (JavaScript)       │ │
│  │  • Syntax Highlighting                │ │
│  │  • Code Completion                    │ │
│  │  • Error Detection                    │ │
│  └───────────────┬───────────────────────┘ │
│                  │                         │
│                  ▼                         │
│  ┌───────────────────────────────────────┐ │
│  │    WebAssembly Runtime (WASM)         │ │
│  │  ┌─────────────────────────────────┐  │ │
│  │  │   Go Runtime + Yaegi            │  │ │
│  │  │   Interpreter                   │  │ │
│  │  └─────────────────────────────────┘  │ │
│  └───────────────┬───────────────────────┘ │
│                  │                         │
│                  ▼                         │
│  ┌───────────────────────────────────────┐ │
│  │         Output Display                │ │
│  │  • Success / Error States             │ │
│  │  • Formatted Output                   │ │
│  └───────────────────────────────────────┘ │
│                                             │
└─────────────────────────────────────────────┘
```

### How It Works

1. **User writes Go code** in the Monaco editor
2. **Code is sent to WASM module** when "Run" is clicked
3. **Yaegi interprets the Go code** within the browser's WASM runtime
4. **Results are captured** and displayed in the output panel
5. **Everything happens locally** — no server communication required

---

## 📁 Project Structure

```
online-go-interpreter/
├── src/
│   ├── index.html          # Main HTML file with Monaco editor integration
│   ├── style.css           # Styling for the application
│   └── build/              # Built WASM files and Go runtime
│       ├── main.wasm       # Compiled Go interpreter
│       └── wasm_exec.js    # Go WASM runtime
├── go-runner/
│   ├── main.go             # Go source for WASM module
│   └── go.mod              # Go module dependencies
├── scripts/
│   └── copy-wasm.js        # Build script to copy WASM runtime files
├── package.json            # npm configuration
└── README.md               # This file
```

---

## 🛠️ Available Scripts

| Command | Description |
|---------|-------------|
| `npm run build` | Compiles the Go code to WebAssembly |
| `npm run serve` | Starts a local development server |
| `npm run dev` | Builds and serves the application |
| `npm run prebuild` | Copies WASM runtime files (runs automatically before build) |

---

## 🎓 Usage Examples

### Hello World

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Working with Variables

```go
package main

import "fmt"

func main() {
    name := "Gopher"
    age := 13
    fmt.Printf("Hello, I'm %s and I'm %d years old!\n", name, age)
}
```

### Functions and Loops

```go
package main

import "fmt"

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
    for i := 0; i < 10; i++ {
        fmt.Printf("fibonacci(%d) = %d\n", i, fibonacci(i))
    }
}
```

---

## 🔧 Technical Details

### Technologies Used

- **Frontend**: HTML5, CSS3, JavaScript (ES6+)
- **Editor**: [Monaco Editor](https://microsoft.github.io/monaco-editor/) v0.44.0
- **UI Components**: [SweetAlert2](https://sweetalert2.github.io/) for modals
- **Go Interpreter**: [Yaegi](https://github.com/traefik/yaegi) — A Go interpreter written in Go
- **Build Tools**: Node.js, npm, cross-env
- **Server**: http-server (for local development)

### Browser Compatibility

- ✅ Chrome/Edge (v57+)
- ✅ Firefox (v52+)
- ✅ Safari (v11+)
- ✅ Opera (v44+)

*Requires WebAssembly support*

---

## 🤝 Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Workflow

1. Make your changes in the appropriate files:
   - UI changes → `src/index.html` or `src/style.css`
   - Interpreter logic → `go-runner/main.go`
   - Build process → `package.json` or `scripts/`

2. Test your changes:
   ```bash
   npm run dev
   ```

3. Ensure everything builds without errors before submitting

---

## 📝 License

This project is licensed under the ISC License. See the repository for details.

---

## 🙏 Acknowledgments

- **[Yaegi](https://github.com/traefik/yaegi)** — The Go interpreter that powers this project
- **[Monaco Editor](https://microsoft.github.io/monaco-editor/)** — The code editor that powers VS Code
- **[The Go Team](https://go.dev/)** — For the amazing Go language and WebAssembly support
- **[SweetAlert2](https://sweetalert2.github.io/)** — For beautiful, responsive UI dialogs

---

## 👨‍💻 Author

**Joshua Beatty**

- Portfolio: [joshbeatty.me](https://joshbeatty.me)
- GitHub: [@Joshua-Beatty](https://github.com/Joshua-Beatty)

---

## 🌟 Star History

If you find this project useful, please consider giving it a ⭐ on GitHub!

---

<div align="center">

**[⬆ Back to Top](#-online-go-interpreter)**

Made with ❤️ and Go

</div>
