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
  root:
    required:
      - README.md
      - framed.yaml
      - go.mod
      - go.sum
      - .gitignore
    children:
      - src:
          required:
            pattern: "*.go"
            files:
              - main.go
          forbidden:
            pattern: "*_test.go" # Disallow tests in /src
      - pipelines:
          required:
            pattern:
              - "*.yml"
              - "*.yaml" # only yaml files allowed
            files:
              - deployment.yaml
              - pr.yaml
      - dockerfiles:
          required:
            pattern:
              - "*.dockerfile"
            minCount: 1 # At least one file has to be there
      - docs:
          required:
            maxCount: 10 # No more than 10 files per dir
            allowChildren: true
            pattern:
              - "*.md"
              - "*.txt"
```

1. **Project Structure Definition**: FRAMED allows you to define the desired structure of your project using a YAML-based configuration file. The configuration specifies the required files and directories that should exist in the project.

2. **Root-level Requirements**: The `root` section defines the files that are required at the root level of the project. These files must be present for the project to be considered valid.

3. **Nested Structure**: The `children` section allows you to define nested directories within the project structure. Each child directory can have its own set of required files and directories.

4. **File Requirements**: You can specify file requirements using the `required` property. It ensures that specific files are present within the designated directory.

5. **File Patterns**: The `pattern` property enables you to define file patterns using glob syntax. This allows for more flexible matching of files based on their extensions or naming conventions.

6. **Forbidden Files**: The `forbidden` property lets you specify file patterns that are not allowed within a directory. This can be useful for enforcing certain naming conventions or excluding specific types of files.

7. **Minimum File Count**: The `minCount` property allows you to set a minimum count for files within a directory. It ensures that a certain number of files must be present in the directory.

8. **Maximum File Count**: The `maxCount` property allows you to set a maximum count for files within a directory. It limits the number of files that can exist within the directory.

9. **Allowing Children**: The `allowChildren` property, when set to true, permits the presence of additional directories within a specified directory. This provides flexibility for organizing files and directories within the project.

By using FRAMED with the provided configuration, you can ensure that your project adheres to a predefined structure, with the required files and directories in place. This helps maintain consistency, reusability, and organization across different projects or within a single project.

## Installation

To install FRAMED, follow these steps:

1. Download the latest release from the FRAMED repository: [link-to-releases]

2. Extract the downloaded archive to a directory of your choice.

3. Add the FRAMED binary to your system's PATH environment variable.

4. Verify the installation by running the following command in your terminal:
framed --version

For detailed installation instructions and alternative installation methods, refer to the [documentation](link-to-installation-guide).

## Usage

### 1. Creating a Project Structure

To create a new project structure using a YAML template, run the following command:

```bash
framed create --template <template-file> --output <destination-path>
```

Replace `<template-file>` with the path to your YAML template file and `<destination-path>` with the directory where you want to generate the project structure.

### 2. Capturing current project structure

To capture the current project structure as a YAML template, run the following command:

```bash
framed capture --output <template-file>
```

### 3. Test Project Structure (CI/CD)

To test the project structure for consistency and compliance with the YAML template, run the following command:

```bash
framed test --template <template-file>
```

For a complete list of available commands and usage examples, refer to the [documentation](link-to-full-docs).

### 4. Visualize Project Structure

To visualize the project structure, run the following command:

```bash
framed visualize --template <template-file>
```