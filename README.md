### README: GitHub Blueprint Searcher

This Go program, when executed, searches for YAML files within a specific GitHub repository and its subdirectories. The program filters files based on a predefined search term (in this case, "blueprint"). Here's how the code works and how you can use it:

### How It Works

1. **GitHub API Integration:**
   - The program uses the GitHub API to fetch content details of a specific path in a GitHub repository.
   - It authenticates the API request using a GitHub personal access token provided as the `accessToken` constant.

2. **Search Functionality:**
   - The `searchFiles` function is a recursive function that explores directories.
   - It checks if a content item is a file, has a YAML extension, and contains the search term in its name.
   - If these conditions are met, the program prints the HTML URL of the YAML file.

3. **Recursive Search:**
   - The program continues to explore subdirectories, making it useful for finding files within a deeply nested directory structure.

### How to Use

1. **Set Up GitHub Access Token:**
   - Replace `"your github token"` in the code with your valid GitHub personal access token.
   - This token is necessary to authenticate requests to the GitHub API and access private repositories if required.

2. **Define Search Parameters:**
   - Modify the `repoOwner`, `repoName`, `searchTerm`, and `targetPath` constants to match your use case.
     - `repoOwner`: GitHub username or organization name.
     - `repoName`: Name of the GitHub repository to search within.
     - `searchTerm`: The term to search for in file names (e.g., "blueprint").
     - `targetPath`: The path within the repository to start the search (e.g., "examples").

3. **Execute the Program:**
   - Run the Go program.
   - The program will print the HTML URLs of YAML files in the specified repository and its subdirectories that match the search criteria.

### Example Usage

```go
go run BlueprintExamples.go
```

### Note
- Ensure your GitHub personal access token has appropriate permissions to access the repository and its contents.
- This code provides a foundation for building more complex GitHub repository search functionality, enabling you to find specific files based on your requirements.
