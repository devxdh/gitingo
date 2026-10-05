# gig

I wanted to learn how Git works under the hood, so I built **gig** (*Git-in-Go*) - a simple Git implementation written in Go.

It isn't a wrapper around the `git` command line tool. It creates standard Git objects (blobs, trees, commits), compresses them with zlib, hashes them with SHA-1, and updates `.git/refs/heads/main` directly.

Since it uses Git's exact file structure and object formats, **you can run gig inside a standard Git repository, or check gig commits using official Git.**

---

## Under the hood

Just Go and its standard library:

- **Go**: For fast execution and single binary builds
- **Cobra** (`spf13/cobra`): For CLI flags and subcommands
- **`crypto/sha1`**: Computes standard 40-character Git object IDs
- **`compress/zlib`**: Reads and writes compressed Git object files in `.git/objects/`

---

## Installing it globally

Options to install `gig` globally across Linux, macOS, or Windows:

### 1. Using Go (`go install`)
If you have Go installed:
```bash
go install github.com/devxdh/gig@latest
```
*(Make sure `$GOPATH/bin` or `~/go/bin` is in your `PATH`.)*

### 2. One-liner shell script (Linux / macOS / WSL)
```bash
curl -fsSL https://raw.githubusercontent.com/devxdh/gig/main/install.sh | bash
```

### 3. Direct binary download
Grab precompiled binaries for Linux, macOS, or Windows directly from [GitHub Releases](https://github.com/devxdh/gig/releases).

---

## Usage

Here is how a basic workflow looks:

```bash
# 1. Initialize a repository (creates .git directory & main branch ref)
gig init my-project
cd my-project

# 2. Add a file to gig's object store (creates a blob object)
echo "Hello world" > main.txt
gig hash-object -w main.txt

# 3. Take a snapshot of the directory (creates a tree object)
gig write-tree

# 4. Create a commit
gig commit -m "Initial commit"

# 5. View commit history
gig log
```

### Branches & Checkout

```bash
# List branches
gig branch

# Create a new branch
gig branch feature

# Switch branches
gig checkout feature
```

### Low-level plumbing commands

Commands for inspecting Git objects:

- `gig cat-file -p <hash>`: Print object contents (blob, tree, or commit)
- `gig cat-file -t <hash>`: Print object type (`blob`, `tree`, or `commit`)
- `gig ls-tree <hash>`: List items in a tree object
- `gig commit-tree <tree-hash> -m "msg" -p <parent-hash>`: Create a commit object manually

---

## Using Gig and official Git together (Hybrid mode)

Because gig writes standard Git objects, you can use it alongside official `git` on the exact same project:

```bash
# Initialize a repo with gig
gig init hybrid-demo
cd hybrid-demo
echo "made with gig" > file.txt
gig commit -m "Commit 1 from gig"

# Now run real Git in the exact same folder - Git reads gig's commits cleanly
git log
git status

# Add another commit using real Git
echo "added with real git" >> file.txt
git commit -am "Commit 2 from real git"

# Check log again with gig
gig log --oneline
```

Both tools read from and write to `.git/objects` and `.git/refs/heads/main`.

---

## Building from source

```bash
git clone https://github.com/devxdh/gig.git
cd gig

make build    # builds binary to ./bin/gig
make install  # installs binary globally to system PATH
make test     # runs unit tests
```

---

## License

[MIT](LICENSE).
