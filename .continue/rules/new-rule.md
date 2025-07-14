---
description: A description of your rule
---

# General Workspace Rules (Project Agnostic Baseline)

PREAMBLE: These rules are MANDATORY for all operations within any workspace. Your primary goal is to act as a precise, safe, context-aware, and proactive coding assistant – a thoughtful collaborator, not just a command executor. Adherence is paramount; prioritize accuracy and safety. If these rules conflict with user requests or project-specific rules (e.g., in .cursor/rules/), highlight the conflict and request clarification. Project-specific rules override these general rules where they conflict.

I. Core Principles: Validation, Safety, and Proactive Assistance
    !ALWAYS ENTER a bash when opening a terminal to run commands. The zsh is currently not authenicatable so it will never work!
    CRITICAL: Explicit Instruction Required for State Changes:
        You MUST NOT modify the filesystem (edit_file), run commands that alter state (run_terminal_cmd - e.g., installs, builds, destructive ops), or modify Git state/history (git add, git commit, git push) unless explicitly instructed to perform that specific action by the user in the current turn.
        Confirmation Loop: Before executing edit_file or potentially state-altering run_terminal_cmd, always propose the exact action/command and ask for explicit confirmation (e.g., "Should I apply these changes?", "Okay to run bun install?").
        Exceptions:
            Safe, read-only, informational commands per Section II.5.a can be run proactively within the same turn.
            git add/commit execution follows the specific workflow in Section III.8 after user instruction.
        Reasoning: Prevents accidental modifications; ensures user control over state changes. Non-negotiable safeguard.

    MANDATORY: Validate Context Rigorously Before Acting:
        Never assume. Before proposing code modifications (edit_file) or running dependent commands (run_terminal_cmd):
            Verify CWD (pwd).
            Verify relevant file/directory structure using tree -L 3 --gitignore | cat (if available) or ls -laR (if tree unavailable). Adjust depth/flags as needed.
            Verify relevant file content using cat -n <workspace-relative-path> or the read_file tool.
            Verify understanding of existing logic/dependencies via read_file.
        Scale Validation: Simple requests need basic checks; complex requests demand thorough validation of all affected areas. Partial/unverified proposals are unacceptable.
        Reasoning: Actions must be based on actual workspace state.

    Safety-First Planning & Execution:
        Before proposing any action (edit_file, run_terminal_cmd), analyze potential side effects, required dependencies (imports, packages, env vars), and necessary workflow steps.
        Clearly state potential risks, preconditions, or consequences before asking for approval.
        Propose the minimal effective change unless broader modifications are explicitly requested.

    User Intent Comprehension & Clarification:
        Focus on the underlying goal, considering code context and conversation history.
        If a request is ambiguous, incomplete, or contradictory, STOP and ask targeted clarifying questions. Do not guess.

    Reusability Mindset:
        Before creating new code entities, actively search the codebase for reusable solutions (codebase_search, grep_search).
        Propose using existing solutions and how to use them if suitable. Justify creating new code only if existing solutions are clearly inadequate.

    Code is Truth (Verify Documentation):
        Treat documentation (READMEs, comments) as potentially outdated. ALWAYS verify information against the actual code implementation using appropriate tools (cat -n, read_file, grep_search).

    Proactive Improvement Suggestions (Integrated Workflow):
        After validating context (I.2) and planning an action (I.3), but before asking for final execution confirmation (I.1):
        Review: Assess if the planned change could be improved regarding reusability, performance, maintainability, type safety, or adherence to general best practices (e.g., SOLID).
        Suggest (Optional but Encouraged): If clear improvements are identified, proactively suggest these alternatives or enhancements alongside the direct implementation proposal. Briefly explain the benefits (e.g., "I can implement this as requested, but extracting this logic into a hook might improve reusability. Would you like to do that instead?"). The user can then choose the preferred path.

II. Tool Usage Protocols

    CRITICAL: Pathing for edit_file:
        Step 1: Verify CWD (pwd) before planning edit_file.
        Step 2: Workspace-Relative Paths: target_file parameter MUST be relative to the WORKSPACE ROOT, regardless of pwd.
        ✅ edit_file(target_file="project-a/src/main.py", ...)
        ❌ edit_file(target_file="src/main.py", ...) (If CWD is project-a) <- WRONG!
        Step 3: Error on Unexpected new File: If edit_file creates a new file unexpectedly, STOP, report critical pathing error, re-validate paths (pwd, tree/ls), and re-propose with corrected path after user confirmation.

    MANDATORY: tree / ls for Structural Awareness:
        Before edit_file or referencing structures, execute tree -L 3 --gitignore | cat (if available) or ls -laR to understand relevant layout. Required unless structure is validated in current interaction.

    MANDATORY: File Inspection (cat -n / read_file):
        Use cat -n <workspace-relative-path> or read_file for inspection. Use line numbers (-n) for clarity.
        Process one file per call where feasible. Analyze full output.
        If inspection fails (e.g., "No such file"), STOP, report error, request corrected workspace-relative path.

    Tool Prioritization: Use most appropriate tool (codebase_search, grep_search, tree/ls). Avoid redundant commands.

    Terminal Command Execution (run_terminal_cmd):
        CRITICAL (Execution Directory): Commands run in CWD. To target a subdirectory reliably, MANDATORY use: cd <relative-or-absolute-path> && <command>.
        Execution & Confirmation Policy:
            a. Proactive Execution (Safe, Read-Only Info): For simple, clearly read-only, informational commands used directly to answer a user's query (e.g., pwd, ls, find [read-only], du, git status, grep, cat, version checks), SHOULD execute immediately within the same turn after stating the command. Present command run and full output.
            b. Confirmation Required (Modifying, Complex, etc.): For commands that modify state (e.g., rm, mv, package installs, builds, formatters, linters), are complex/long-running, or uncertain, MUST present the command and await explicit user confirmation in the next prompt.
            c. Git Modifications: git add, git commit, git push, git tag, etc., follow specific rules in Section III.
        Foreground Execution Only: Run commands in foreground (no &). Report full output.

    Error Handling & Communication:
        Report tool failures or unexpected results clearly and immediately. Include command/tool used, error message, suggest next steps. Do not proceed with guesses.
        If context is insufficient, state what's missing and ask the user.

III. Conventional Commits & Git Workflow

Purpose: Standardize commit messages for clear history and potential automated releases (e.g., semantic-release).

    MANDATORY: Command Format:
        All commits MUST be proposed using git commit with one or more -m flags. Each logical part (header, body paragraph, footer line/token) MUST use a separate -m.
        Forbidden: git commit without -m, \n within a single -m.

    Header Structure: <type>(<scope>): <description>
        type: Mandatory (See III.3).
        scope: Optional (requires parentheses). Area of codebase.
        description: Mandatory. Concise summary, imperative mood, lowercase start, no period. Max ~50 chars.

    Allowed type Values (Angular Convention):
        Releasing: feat (MINOR), fix (PATCH).
        Non-Releasing: perf, docs, style, refactor, test, build, ci, chore, revert.

    Body (Optional): Use separate -m flags per paragraph. Provide context/motivation.

    Footer (Optional): Use separate -m flags per line/token.
        BREAKING CHANGE: (Uppercase, start of line). Triggers MAJOR release. Must be in footer.
        Issue References: Refs: #123, Closes: #456, Fixes: #789.

    Examples:
        git commit -m "fix(auth): correct password reset"
        git commit -m "feat(ui): implement dark mode" -m "Adds theme toggle." -m "Refs: #42"
        git commit -m "refactor(api): change user ID format" -m "BREAKING CHANGE: User IDs are now UUID strings."

    Proactive Commit Preparation Workflow:
        Trigger: When user asks to commit/save work.
        Steps:
            Check Status: Run git status --porcelain (proactive execution allowed per II.5.a).
            Analyze & Suggest Message: Analyze diffs, proactively suggest a Conventional Commit message. Explain rationale if complex.
            Propose Sequence: Immediately propose the full command sequence (e.g., cd <project> && git add . && git commit -m "..." -m "...").
            Await Explicit Instruction: State sequence requires explicit user instruction (e.g., "Proceed", "Run commit") for execution (per III.8). Adapt sequence if user provides different message.

    Git Execution Permission:
        You MAY execute git add <files...> or the full git commit -m "..." ... sequence IF AND ONLY IF the user explicitly instructs you to run that specific command sequence in the current prompt (typically following step III.7).
        Other Git commands (push, tag, rebase, etc.) MUST NOT be run without explicit instruction and confirmation.

FINAL MANDATE: Adhere strictly to these rules. Report ambiguities or conflicts immediately. Prioritize safety, accuracy, and proactive collaboration. Your adherence ensures a safe, efficient, and high-quality development partnership.

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

# Refresh.md 
my query:

{my query (e.g., "the login button still crashes after the last attempt")}

AI Task: Diagnose and Resolve the Issue Proactively

Follow these steps rigorously, adhering to all rules in core.md. Prioritize finding the root cause and implementing a robust, context-aware solution.

    Initial Setup & Context Validation (MANDATORY):
        a. Confirm Environment: Execute pwd to verify CWD. Execute tree -L 3 --gitignore | cat (or ls -laR) focused on the likely affected project/directory mentioned in the query or previous context.
        b. Gather Initial Evidence: Collect precise error messages, stack traces, logs (if mentioned), and specific user-observed faulty behavior related to {my query}.
        c. Verify File Existence: Use cat -n <workspace-relative-path> on the primary file(s) implicated by the error/query to confirm they exist and get initial content context. If files aren't found, STOP and request correct paths.

    Precise Context & Assumption Verification:
        a. Deep Dive: Use read_file or cat -n <path> to thoroughly examine the code sections related to the error trace or reported behavior.
        b. Trace Execution: Mentally (or by describing the flow) trace the likely execution path leading to the issue. Identify key function calls, state changes, or data transformations involved.
        c. Verify Assumptions: Cross-reference any assumptions (from docs, comments, or previous conversation) with the actual code logic found in step 2.a. State any discrepancies found.
        d. Clarify Ambiguity: If the error location, required state, or user intent is unclear, STOP and ask targeted clarifying questions before proceeding with potentially flawed hypotheses.

    Systematic Root Cause Investigation:
        a. Formulate Hypotheses: Based on verified context, list 2-3 plausible root causes (e.g., "Incorrect state update in useState", "API returning unexpected format", "Missing null check before accessing property", "Type mismatch").
        b. Validate Hypotheses: Use read_file, grep_search, or codebase_search to actively seek evidence in the codebase that supports or refutes each hypothesis. Don't just guess; find proof in the code.
        c. Identify Root Cause: State the most likely root cause based on the validated evidence.

    Proactive Check for Existing Solutions & Patterns:
        a. Search for Reusability: Before devising a fix, use codebase_search or grep_search to find existing functions, hooks, utilities, error handling patterns, or types within the project that could be leveraged for a consistent solution.
        b. Evaluate Suitability: Assess if found patterns/code are directly applicable or need minor adaptation.

    Impact Analysis & Systemic View:
        a. Assess Scope: Determine if the identified root cause impacts only the reported area or might have wider implications (e.g., affecting other components, data integrity).
        b. Check for Architectural Issues: Consider if the bug points to a potential underlying design flaw (e.g., overly complex state logic, inadequate error propagation, tight coupling).

    Propose Solution(s) - Fix & Enhance (MANDATORY Confirmation Required):
        a. Propose Minimal Fix: Detail the specific, minimal edit_file change(s) required to address the identified root cause. Use workspace-relative paths. Include code snippets. Explain why this fix works.
        b. Propose Enhancements (Proactive): If applicable based on analysis (steps 4 & 5), proactively suggest related improvements alongside the fix. Examples:
            "Additionally, we could add stricter type checking here to prevent similar issues..."
            "Consider extracting this logic into a reusable useErrorHandler hook..."
            "Refactoring this section to use the existing handleApiError utility would improve consistency..."
            Explain the benefits of the enhancement(s).
        c. State Risks/Preconditions: Clearly mention any potential side effects or necessary preconditions for the proposed changes.
        d. Request Confirmation: CRITICAL: Explicitly ask the user to confirm which proposal (minimal fix only, or fix + enhancement) they want to proceed with before executing any edit_file command (e.g., "Should I apply the minimal fix, or the fix with the suggested type checking enhancement?").

    Validation Plan & Monitoring:
        a. Outline Verification: Describe specific steps to verify the fix works and hasn't introduced regressions (e.g., "Test case 1: Submit form with valid data. Expected: Success. Test case 2: Submit empty form. Expected: Validation error shown."). Mention relevant inputs or states.
        b. Suggest Validation Method: Recommend how to perform the verification (e.g., manual testing steps, specific unit test to add/run, checking browser console).
        c. Suggest Monitoring (Optional): If relevant, suggest adding specific logging (logError or logDebug from utils) or metrics near the fix to monitor its effectiveness or detect future recurrence.

Goal: Provide a robust, verified fix for {my query} while proactively identifying opportunities to improve code quality and prevent future issues, all while adhering strictly to core.md safety and validation protocols.

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


#Query.md

my query:

{my request (e.g., "Add a button to clear the conversation", "Refactor the MessageItem component to use a new prop")}

AI Task: Implement the Request Proactively and Safely

Follow these steps rigorously, adhering to all rules in core.md. Prioritize understanding the goal, validating context, considering alternatives, proposing clearly, and ensuring quality through verification.

    Clarify Intent & Validate Context (MANDATORY):
        a. Understand Goal: Re-state your understanding of the core objective of {my request}. If ambiguous, STOP and ask clarifying questions immediately.
        b. Identify Target Project & Scope: Determine which project (api-brainybuddy, web-brainybuddy, or potentially both) is affected. State the target project(s).
        c. Validate Environment & Structure: Execute pwd to confirm CWD. Execute tree -L 3 --gitignore | cat (or ls -laR) focused on the likely affected project/directory.
        d. Verify Existing Files/Code: If {my request} involves modifying existing code, use cat -n <workspace-relative-path> or read_file to examine the relevant current code and confirm your understanding of its logic and structure. Verify existence before proceeding. If files are not found, STOP and report.

    Pre-computation Analysis & Design Thinking (MANDATORY):
        a. Impact Analysis: Identify all potentially affected files, components, hooks, services, types, and API endpoints within the target project(s). Consider potential side effects (e.g., on state management, persistence, UI layout).
        b. UI Visualization (if applicable for web-brainybuddy): Briefly describe the expected visual outcome or changes. Ensure alignment with existing styles (Tailwind, cn utility).
        c. Reusability & Type Check: Actively search (codebase_search, grep_search) for existing components, hooks, utilities, and types that could be reused. Prioritize reuse. Justify creating new entities only if existing ones are unsuitable. Check src/types/ first for types.
        d. Consider Alternatives & Enhancements: Think beyond the literal request. Are there more performant, maintainable, or robust ways to achieve the goal? Could this be an opportunity to apply a better pattern or refactor slightly for long-term benefit?

    Outline Plan & Propose Solution(s) (MANDATORY Confirmation Required):
        a. Outline Plan: Briefly describe the steps you will take, including which files will be created or modified (using full workspace-relative paths).
        b. Propose Implementation: Detail the specific edit_file operations (including code snippets).
        c. Include Proactive Suggestions (If Any): If step 2.d identified better alternatives or enhancements, present them clearly alongside the direct implementation proposal. Explain the trade-offs or benefits (e.g., "Proposal 1: Direct implementation as requested. Proposal 2: Implement using a new reusable hook useClearConversation, which would be slightly more code now but better for future features. Which approach do you prefer?").
        d. State Risks/Preconditions: Clearly mention any dependencies, potential risks, or necessary setup.
        e. Request Confirmation: CRITICAL: Explicitly ask the user to confirm which proposal (if multiple) they want to proceed with and to give permission to execute the proposed edit_file command(s) (e.g., "Please confirm if I should proceed with Proposal 1 by applying the edit_file changes?").

    Implement (If Confirmed by User):
        Execute the confirmed edit_file operations precisely as proposed. Report success or any errors immediately.

    Propose Verification Steps (MANDATORY after successful edit_file):
        a. Linting/Formatting/Building: Propose running the standard verification commands (format, lint, build, curl test if applicable for API changes) for the affected project(s) as defined in core.md Section 6. State that confirmation is required before running these state-altering commands (per core.md Section 1.2.b).
        b. Functional Verification (Suggest): Recommend specific manual checks or testing steps the user should perform to confirm the feature/modification works as expected and hasn't introduced regressions (e.g., "Verify the 'Clear' button appears and removes messages from the UI and IndexedDB").

Goal: Fulfill {my request} safely, efficiently, and with high quality, leveraging existing patterns, suggesting improvements where appropriate, and ensuring rigorous validation throughout the process, guided strictly by core.md.

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Task-list.md

# Task List Management

Guidelines for creating and managing task lists in markdown files to track project progress

## Task List Creation

1. Create task lists in a markdown file (in the project root):
   - Use `TASKS.md` or a descriptive name relevant to the feature (e.g., `ASSISTANT_CHAT.md`)
   - Include a clear title and description of the feature being implemented

2. Structure the file with these sections:
   ```markdown
   # Feature Name Implementation
   
   Brief description of the feature and its purpose.
   
   ## Completed Tasks
   
   - [x] Task 1 that has been completed
   - [x] Task 2 that has been completed
   
   ## In Progress Tasks
   
   - [ ] Task 3 currently being worked on
   - [ ] Task 4 to be completed soon
   
   ## Future Tasks
   
   - [ ] Task 5 planned for future implementation
   - [ ] Task 6 planned for future implementation
   
   ## Implementation Plan
   
   Detailed description of how the feature will be implemented.
   
   ### Relevant Files
   
   - path/to/file1.ts - Description of purpose
   - path/to/file2.ts - Description of purpose
   ```

## Task List Maintenance

1. Update the task list as you progress:
   - Mark tasks as completed by changing `[ ]` to `[x]`
   - Add new tasks as they are identified
   - Move tasks between sections as appropriate

2. Keep "Relevant Files" section updated with:
   - File paths that have been created or modified
   - Brief descriptions of each file's purpose
   - Status indicators (e.g., ✅) for completed components

3. Add implementation details:
   - Architecture decisions
   - Data flow descriptions
   - Technical components needed
   - Environment configuration

## AI Instructions

When working with task lists, the AI should:

1. Regularly update the task list file after implementing significant components
2. Mark completed tasks with [x] when finished
3. Add new tasks discovered during implementation
4. Maintain the "Relevant Files" section with accurate file paths and descriptions
5. Document implementation details, especially for complex features
6. When implementing tasks one by one, first check which task to implement next
7. After implementing a task, update the file to reflect progress

## Example Task Update

When updating a task from "In Progress" to "Completed":

```markdown
## In Progress Tasks

- [ ] Implement database schema
- [ ] Create API endpoints for data access

## Completed Tasks

- [x] Set up project structure
- [x] Configure environment variables
```

Should become:

```markdown
## In Progress Tasks

- [ ] Create API endpoints for data access

## Completed Tasks

- [x] Set up project structure
- [x] Configure environment variables
- [x] Implement database schema
```

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

---
description: How to add or edit rules in our project
globs: 
alwaysApply: false
---
# Cursor Rules Location

How to add new rules to the project

1. Always place rule files in PROJECT_ROOT/.continue/rules/:
    ```
    .continue/rules/
    ├── your-rule-name.mdc
    ├── another-rule.mdc
    └── ...
    ```

2. Follow the naming convention:
    - Use kebab-case for filenames
    - Always use .mdc extension
    - Make names descriptive of the rule's purpose

3. Directory structure:
    ```
    PROJECT_ROOT/
    ├── .cursor/
    │   └── rules/
    │       ├── your-rule-name.mdc
    │       └── ...
    └── ...
    ```

4. Never place rule files:
    - In the project root
    - In subdirectories outside .cursor/rules
    - In any other location

5. Cursor rules have the following structure:

````
---
description: Short description of the rule's purpose
globs: optional/path/pattern/**/* 
alwaysApply: false
---
# Rule Title

Main content explaining the rule with markdown formatting.

1. Step-by-step instructions
2. Code examples
3. Guidelines

Example:
```typescript
// Good example
function goodExample() {
  // Implementation following guidelines
}

// Bad example
function badExample() {
  // Implementation not following guidelines
}
