# Side-by-Side Git Diff TUI

A Terminal User Interface (TUI) application written in Go to display a side-by-side git diff with synchronized scrolling.

## Features

- **Side-by-side diff display**: View old and new versions simultaneously
- **Synchronized scrolling**: Both panes scroll together for easy comparison
- **Line numbers**: Each line shows its position in the respective file
- **Syntax highlighting**: Added lines in green, removed lines in red
- **Mouse support**: Click and scroll with your mouse
- **Keyboard navigation**: Use arrow keys, j/k, or mouse wheel

## How to Run

1.  **Prerequisites:** Ensure you have [Go](https://go.dev/doc/install) and [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) installed on your system.
2.  **Navigate to Project:** Open your terminal in the project's root directory.
3.  **Install Dependencies:** If you haven't already, run the following command to download the necessary Go modules:
    ```sh
    go mod tidy
    ```
4.  **Run the Application:** Execute the following command to start the TUI. The application shows the output of `git diff`, so make sure you have staged or unstaged changes in your local repository to see a result.
    ```sh
    go run main.go
    ```

## Controls

-   **Scroll:** Use the arrow keys (`↑`/`↓`), `j`/`k`, or your mouse wheel to scroll through the diff. Both panes will scroll simultaneously.
-   **Quit:** Press `q`, `esc`, or `ctrl+c` to exit the application.