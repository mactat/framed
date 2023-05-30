# FRAMED (Files and Directories Reusability, Architecture, and Management)

FRAMED is a powerful CLI tool written in Go that simplifies the organization and management of files and directories in a reusable and architectural manner. It provides YAML templates for defining project structures and enables workflows based on those.

To always be in sync with the YAML template, FRAMED provides a built-in test command that can be used in CI/CD pipelines to verify the project structure.

## Features

- **YAML Templates**: FRAMED uses YAML templates to define the entire project structure.

- **Always in Sync**: FRAMED provides a built-in test command that can be used in CI/CD pipelines to verify the project structure and ensure that it is always in sync with the YAML template.

- **Consistency Across Projects**: FRAMED offers a consistent way of organizing files and directories across different projects.

## Example configuration

To get started with FRAMED, you can use the following example:

```yaml
# FRAMED Configuration
name: framed

structure:
name: framed

structure:
  name: root
  maxDepth: 5 # Disallow dirs deeper than 5
  files:
    - README.md
    - framed.yaml
    - main.go
    - go.mod
    - go.sum
    - .gitignore
  dirs:
    - name: cmd
      allowedPatterns:
        - ".go"
      forbiddenPatterns:
        - "_test.go" # Disallow tests in /src
    - name: pipelines
      maxCount: 2
      allowedPatterns:
        - ".yml"
        - ".yaml" # only yaml files allowed
      files:
        - pr.yaml
    - name: dockerfiles
      minCount: 1 # At least one file has to be there
      allowChildren: false # Allow subdirectories creation, default true
      allowedPatterns:
        - ".dockerfile"
    - name: docs
      maxCount: 10 # No more than 10 files per dir
      allowedPatterns:
        - ".md"
        - ".txt"
      dirs:
        - name: design
    - name: examples

    - name: other
      template: other.yaml # Use another template for this dir
    - name: another
      template: https://github.com/mactat/framed/blob/master/framed.yaml # Share templates between projects

```

### Project Structure Definition

FRAMED allows you to define the desired structure of your project using a YAML-based configuration file. The configuration specifies the required files and directories that should exist in the project.

### Root-level Requirements

The structure section defines the files that are required at the root level of the project. These files must be present for the project to be considered valid.

### Nested Structure

The dirs section allows you to define nested directories within the project structure. Each subdirectory can have its own set of required files and directories.

### File Requirements

You can specify file requirements using the files property. It ensures that specific files are present within the designated directory.

### File Patterns

The allowedPatterns property enables you to define file patterns using glob syntax. This allows for more flexible matching of files based on their extensions or naming conventions.

### Forbidden Files

The forbiddenPatterns property lets you specify file patterns that are not allowed within a directory. This can be useful for enforcing certain naming conventions or excluding specific types of files.

### Minimum File Count

The minCount property allows you to set a minimum count for files within a directory. It ensures that a certain number of files must be present in the directory.

### Maximum File Count

The maxCount property allows you to set a maximum count for files within a directory. It limits the number of files that can exist within the directory.

### Allowing Children

The allowChildren property, when set to true, permits the presence of additional directories within a specified directory. This provides flexibility for organizing files and directories within the project.

## Installation

To install FRAMED, follow these steps:

1. Download the latest release from the FRAMED repository: [link-to-releases]

2. Extract the downloaded archive to a directory of your choice.

3. Add the FRAMED binary to your system's PATH environment variable.

4. Verify the installation by running the following command in your terminal:
framed --version

For detailed installation instructions and alternative installation methods, refer to the [documentation](link-to-installation-guide).

## Usage

**Note**: The following commands assume that you have already installed FRAMED and added it to your system's PATH environment variable.

**Note**: By default template file is `framed.yaml`. You can specify a different template file using the `--template` flag f.e `--template path/to/my-template.yaml`.

### 1. Creating a Project Structure

To create a new project structure using a YAML template, run the following command:

```bash
framed create
```

If you also want to create required files, run the following command:

```bash
framed create --files
```

### 2. Capturing current project structure

To capture the current project structure as a YAML template, run the following command:

```bash
framed capture --output <template-file>
```

### 3. Test Project Structure (CI/CD)

To test the project structure for consistency and compliance with the YAML template, run the following command:

```bash
framed verify
```

For a complete list of available commands and usage examples, refer to the [documentation](link-to-full-docs).

### 4. Visualize Project Structure

To visualize the project structure, run the following command:

```bash
framed visualize
```

## Demo

![Demo](./docs/static/demo.gif)