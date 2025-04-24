# golang-project-template
Starting point for a Go project. Ready to use with Docker, Taskfile and Lefthook. You don't need to worry about downloading golang, as it is already included in the Docker image.

## Requirements
- [docker](https://www.docker.com/get-started/)
- [task command](https://taskfile.dev/installation/) (for executing commands easier)
- [Lefthook](https://github.com/evilmartians/lefthook) (optional: if you want to use pre-commit hooks for github)
- [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm) (optional: if you want to use Lefthook)

## Initialization
If you want to use Lefthook, run the following command:
```
npx lefthook install
```

## Project Commands
### Run the project (Start containers)
```
task
```

### Stop the project (Stop containers)
```
task down
```

### Run golangci-lint
```
task lint
```

### Delete volumes
```
task destroy
```