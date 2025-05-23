# Quiver Compiler

The compiler works as follows:

- Is install on your project root folder.
- Lets you compile almost anything!

## Publishing a Go package

Publishing a Go package, or more accurately, a Go module, involves making your code available for other developers to use. Here's a breakdown of the process:

**4. Tag a Version**

*   To allow others to depend on a specific version of your module, you need to tag your releases. Go modules use semantic versioning (e.g., `v1.0.0`, `v1.0.1`, `v0.1.0`).
*   Use Git to tag your commit:
    ```bash
    git tag v1.0.0
    ```
*   Push the tag to your repository:
    ```bash
    git push origin v1.0.0
    ```
    Or, to push all tags:
    ```bash
    git push --tags
    ```
    Some sources suggest `git push --follow-tags origin` as a way to push the current branch along with its tags.