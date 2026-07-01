```markdown
# spoony Development Patterns

> Auto-generated skill from repository analysis

## Overview
This skill introduces the core development patterns and conventions used in the `spoony` JavaScript codebase. It covers file naming, import/export styles, commit message practices, and testing patterns. While no specific framework or automated workflows are detected, this guide will help you quickly align with the project's standards and streamline your contributions.

## Coding Conventions

### File Naming
- Use **camelCase** for file names.
  - Example: `myUtility.js`, `dataFetcher.js`

### Imports
- Use **relative import paths**.
  - Example:
    ```javascript
    import fetchData from './fetchData';
    import { parseUser } from '../utils/parseUser';
    ```

### Exports
- Both **default** and **named exports** are used.
  - Default export example:
    ```javascript
    export default function fetchData() { ... }
    ```
  - Named export example:
    ```javascript
    export function parseUser(data) { ... }
    ```

### Commit Messages
- **Freeform** style, no enforced prefixes.
- Average commit message length: ~72 characters.
  - Example:  
    ```
    Add utility to handle API response parsing
    ```

## Workflows

### Adding a New Module
**Trigger:** When you need to add new functionality as a separate module  
**Command:** `/add-module`

1. Create a new file using camelCase naming (e.g., `myFeature.js`).
2. Implement your logic using JavaScript.
3. Use relative imports for any dependencies.
4. Export your module (default or named as appropriate).
5. Write a corresponding test file named `myFeature.test.js`.
6. Commit your changes with a clear, descriptive message.

### Writing and Running Tests
**Trigger:** When you need to verify code correctness  
**Command:** `/run-tests`

1. Create a test file named with the pattern `*.test.js` (e.g., `myFeature.test.js`).
2. Write your test cases (testing framework is not specified; follow project precedent).
3. Run tests using the project's preferred method (consult project documentation if needed).
4. Review test results and fix any failing cases.

## Testing Patterns

- Test files follow the pattern: `*.test.js`
- Testing framework is **unknown**; check existing test files for style and structure.
- Place test files alongside the modules they test or in a dedicated test directory.
- Example test file name: `fetchData.test.js`

## Commands
| Command       | Purpose                                      |
|---------------|----------------------------------------------|
| /add-module   | Scaffold and add a new module                |
| /run-tests    | Run all test files matching `*.test.js`      |
```
