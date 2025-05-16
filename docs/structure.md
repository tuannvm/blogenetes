## Best Practices for Structuring a GitHub Repository for Multiple Dagger Modules

Organizing a GitHub repository to host multiple Dagger modules—especially for a Go project—requires balancing clarity, modularity, and maintainability. Below is a detailed guide on how to structure such a repository, drawing from Dagger’s official documentation and real-world monorepo examples.

---

### 1. **Monorepo vs. Single-Module Repo**

Dagger fully supports both single-module repositories and monorepos containing multiple modules. For your use case—centralizing all your Dagger modules in one repo—a monorepo is ideal. Dagger is agnostic to repository layout, so you can organize modules as you see fit, and each module can have its own dependencies and configuration.

---

### 2. **Recommended Directory Structure**

A typical monorepo for Dagger modules in Go might look like this:

```
dagger-modules-repo/
├── go.mod
├── go.sum
├── README.md
├── .github/
│   └── workflows/
├── module-a/
│   ├── dagger.json
│   ├── main.go
│   ├── internal/
│   └── ... (other Go files)
├── module-b/
│   ├── dagger.json
│   ├── main.go
│   ├── internal/
│   └── ...
├── shared/
│   └── ... (shared Go code, utilities, etc.)
└── docs/
    └── ...
```

**Key Points:**
- Each module (`module-a`, `module-b`, etc.) is a self-contained Dagger module with its own `dagger.json` and source code.
- Shared code or utilities can live in a `shared/` directory at the repo root, which modules can import as needed.
- The root `go.mod` can be used for dependency management across all modules, or you can have separate `go.mod` files per module if you want stricter isolation.
- `.github/workflows/` can contain CI/CD pipelines for building, testing, and publishing your modules.

---

### 3. **Where Should Dagger Functions Live?**

You do **not** need to keep all Dagger functions in a `.dagger` directory. In fact, Dagger recommends placing module source code in a dedicated subdirectory for each module. The `.dagger` directory is not a requirement; it’s just a convention some projects use. You can name your module directories however you like, and keep them at the root or in a `modules/` folder for further organization.

---

### 4. **Module Initialization and Configuration**

- Each module directory should contain its own `dagger.json` file, which configures the module.
- The source code for the module (e.g., `main.go`, `internal/`, etc.) lives alongside `dagger.json`.
- If you want to keep functions in a separate place (e.g., a `shared/` directory), you can import them into your module’s main code as you would with any Go package.

---

### 5. **Accessing Shared Code and Dependencies**

If modules need to share code, place common packages in a `shared/` directory and import them using Go’s import path conventions. If you use a single `go.mod` at the repo root, all modules can reference shared code easily. If you use per-module `go.mod` files, you may need to use Go workspaces or replace directives.

---

### 6. **Example: Two Go Modules with Shared Code**

```
dagger-modules-repo/
├── go.mod
├── module-foo/
│   ├── dagger.json
│   ├── main.go
│   └── internal/
├── module-bar/
│   ├── dagger.json
│   ├── main.go
│   └── internal/
└── shared/
    └── utils.go
```
- `module-foo/main.go` and `module-bar/main.go` can both import `"github.com/your-org/dagger-modules-repo/shared"`.

---

### 7. **Advanced: Contextual Modules and Filtering**

If a module needs access to files outside its directory (e.g., for monorepo builds), you can configure the `source` field in `dagger.json` to point to the relevant context directory. Dagger is evolving to make contextual modules (modules that operate within a larger project context) a first-class concept, so you can expect even better support for this pattern in the future.

---

### 8. **Summary of Best Practices**

- **One directory per module**: Each with its own `dagger.json` and source code.
- **No need for a `.dagger` directory**: Use clear, descriptive names for module directories.
- **Shared code**: Place in a `shared/` directory at the repo root.
- **Root-level `go.mod`**: For easier dependency management and code sharing.
- **Documentation and CI**: Keep docs and workflows at the repo root for visibility and maintainability.

---

### 9. **References and Further Reading**

- [Dagger Module Structure Documentation]
- [Reusable Modules and Monorepo Support]
- [Discussion on Contextual Modules]
- [Example Monorepo Structure]

---

**In summary:**  
Structure your Dagger modules repo as a monorepo, with each module in its own directory (not necessarily `.dagger`), shared code in a common directory, and a root `go.mod` for dependency management. This approach is scalable, maintainable, and aligns with Dagger’s best practices for multi-module projects.

---
: https://docs.dagger.io/features/modules  : https://docs.dagger.io/api/module-structure  : https://github.com/dagger/dagger/issues/6291  : https://github.com/dagger/dagger/issues/7199
