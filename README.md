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

1. **Execute the Program:**
   - Run the Go program.
   - The program will open UI and have fields.
     
2. **Set Up GitHub Access Token:**
   - Add `"your github token"` in the UI with your valid GitHub personal access token.
   - This token is necessary to authenticate requests to the GitHub API and access private repositories if required.

3. **Define Search Parameters:**
   - Add the `searchTerm` constants to match your use case.
     - `searchTerm`: The term to search for in file names (e.g., "MongoD").
     
### Example Usage

```go
go run Kansearch.go
```

### Note
- Ensure your GitHub personal access token has appropriate permissions to access the repository and its contents.
- This code provides a foundation for building more complex GitHub repository search functionality, enabling you to find specific files based on your requirements.

## Screenshot

![image](https://github.com/mritsurgeon/Kanister_BluePrint_Search/assets/59644778/15e159df-cb5d-43a7-8477-870be7b07ebb)

