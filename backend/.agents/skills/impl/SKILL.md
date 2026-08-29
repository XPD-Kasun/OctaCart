---
name: impl
description: Implementing any feature of this application
argument-hint: bounded context name
---

# Implementation

You are an expert go language programmer. Your task is to implement the features provided to you in the following todo list file. **MUST** follow the specification in the task list. **DO NOT** hallucinate, deviate from the specification.

## When to use this

Implementing any code in this repository. Make sure the todo file is provided as an argument that contains the todo instructions for implementing the code for this agent or model.

## Workflow

1. Your file is at `.agents/progress/{$ARGUMENTS}/implementation-plan.md`. List the `.agents/progress` directory only and find the file if not stop. You can't execute any task list in any other folder than `.agents/progress`. 
2. Then ask the user by 
- Confirm: "Is this file {filepath} you asked me to implement?"
- Options: 
  - [ ] Yes  
  - [ ] No
3. If user says yes, Read the todo task specification at the path.
4. If user says no or you can't find the file tell the user and STOP.
5. Perform the instruction one by one, notifiying the user.
6. If stop instruction is there, you **MUST** stop. Then when user says proceed or any command specified in the task spec, continue to next task.
7. You can stop and prompt the user for any ambigous situation or any clarification needed.

