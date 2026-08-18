# gitingo

I wanted to learn how Git works under the hood, so I built **gitingo** - a simple Git implementation written in Go.

It isn't a wrapper around the `git` command line tool. It creates standard Git objects (blobs, trees, commits), compresses them with zlib, hashes them with SHA-1, and updates `.git/refs/heads/main` directly.

Since it uses Git's exact file structure and object formats, **you can run gitingo inside a standard Git repository, or check gitingo commits using official Git.**

---

## Under the hood

Just Go and its standard library:

- **Go**: For fast execution and single binary builds
- **Cobra** (`spf13/cobra`): For CLI flags and subcommands
- **`crypto/sha1`**: Computes standard 40-character Git object IDs
- **`compress/zlib`**: Reads and writes compressed Git object files in `.git/objects/`

---

## Installing it globally

Options to install `gitingo` globally across Linux, macOS, or Windows:

### 1. Using Go (`go install`)
If you have Go installed:
```bash
go install github.com/devxdh/gitingo@latest
```
*(Make sure `$GOPATH/bin` or `~/go/bin` is in your `PATH`.)*

### 2. One-liner shell script (Linux / macOS / WSL)
```bash
curl -fsSL https://raw.githubusercontent.com/devxdh/gitingo/main/install.sh | bash
```

### 3. Direct binary download
Grab precompiled binaries for Linux, macOS, or Windows directly from [GitHub Releases](https://github.com/devxdh/gitingo/releases).

---

## Usage

Here is how a basic workflow looks:

```bash
# 1. Initialize a repository (creates .git directory & main branch ref)
gitingo init my-project
cd my-project

# 2. Add a file to gitingo's object store (creates a blob object)
echo "Hello world" > main.txt
gitingo hash-object -w main.txt

# 3. Take a snapshot of the directory (creates a tree object)
gitingo write-tree

# 4. Create a commit
gitingo commit -m "Initial commit"

# 5. View commit history
gitingo log
```

### Low-level plumbing commands

Commands for inspecting Git objects:

- `gitingo cat-file -p <hash>`: Print object contents (blob, tree, or commit)
- `gitingo cat-file -t <hash>`: Print object type (`blob`, `tree`, or `commit`)
- `gitingo ls-tree <hash>`: List items in a tree object
- `gitingo commit-tree <tree-hash> -m "msg" -p <parent-hash>`: Create a commit object manually

---

## Using Gitingo and official Git together (Hybrid mode)

Because gitingo writes standard Git objects, you can use it alongside official `git` on the exact same project:

```bash
# Initialize a repo with gitingo
gitingo init hybrid-demo
cd hybrid-demo
echo "made with gitingo" > file.txt
gitingo commit -m "Commit 1 from gitingo"

# Now run real Git in the exact same folder - Git reads Gitingo's commits cleanly
git log
git status

# Add another commit using real Git
echo "added with real git" >> file.txt
git commit -am "Commit 2 from real git"

# Check log again with gitingo
gitingo log --oneline
```

Both tools read from and write to `.git/objects` and `.git/refs/heads/main`.

---

## Building from source

```bash
git clone https://github.com/devxdh/gitingo.git
cd gitingo

make build    # builds binary to ./bin/gitingo
make install  # installs binary globally to system PATH
make test     # runs unit tests
```

---

## License

[MIT](LICENSE).
