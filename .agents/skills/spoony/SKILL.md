```markdown
# spoony Development Patterns

> Auto-generated skill from repository analysis

## Overview
This skill teaches the core development patterns and conventions used in the `spoony` JavaScript repository. You'll learn about file naming, import/export styles, commit message habits, and how to structure and run tests. This guide is ideal for contributors aiming to maintain consistency and quality in the codebase.

## Coding Conventions

### File Naming
- Use **camelCase** for all file names.
  - Example: `myUtilityFile.js`

### Import Style
- Use **relative imports** for modules within the project.
  - Example:
    ```javascript
    import helper from './utils/helper';
    ```

### Export Style
- **Mixed exports** are used (both default and named).
  - Example:
    ```javascript
    // Default export
    export default function main() { ... }

    // Named export
    export function helper() { ... }
    ```

### Commit Messages
- **Freeform** commit messages, sometimes with prefixes.
- Average length: ~87 characters.
  - Example:
    ```
    Add support for new API endpoint and update related documentation
    ```

## Workflows

### Adding a New Feature
**Trigger:** When implementing a new feature or functionality  
**Command:** `/add-feature`

1. Create a new file using camelCase naming.
2. Write your feature code, using relative imports for dependencies.
3. Export your feature using default or named exports as appropriate.
4. Add or update corresponding test files (`*.test.*`).
5. Commit your changes with a clear, descriptive message.
6. Open a pull request for review.

### Fixing a Bug
**Trigger:** When resolving a bug or issue  
**Command:** `/fix-bug`

1. Identify the affected file(s) and make necessary code changes.
2. Ensure all imports remain relative and follow export conventions.
3. Update or add test cases to cover the bug fix.
4. Commit your changes with a descriptive message.
5. Open a pull request referencing the issue (if applicable).

### Writing and Running Tests
**Trigger:** When adding or updating tests  
**Command:** `/run-tests`

1. Create or update test files following the `*.test.*` pattern.
2. Write test cases to cover new or changed functionality.
3. Use the project's preferred (unknown) test runner to execute tests.
4. Review test results and fix any failures.

## Testing Patterns

- Test files use the `*.test.*` naming convention.
  - Example: `mathUtils.test.js`
- The test framework is **unknown**, so refer to existing test files for structure and assertions.
- Place test files alongside or near the modules they test.

## Commands
| Command       | Purpose                                 |
|---------------|-----------------------------------------|
| /add-feature  | Start the workflow for adding a feature |
| /fix-bug      | Start the workflow for fixing a bug     |
| /run-tests    | Run all test suites                     |
```