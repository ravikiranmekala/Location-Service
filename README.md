# Project Title

## Description

This project involves creating REST APIs using the Go programming language and the Gin web framework. 
The code can be run in two different modes:

1. **Development Mode:**
   - Hot code reloading is available.
   - Uses Dockerfile1 and docker-compose1.yml.
   - To run, execute the following command:
     ```bash
     docker-compose -f docker-compose1.yml up
     ```

2. **Deployment Mode:**
   - No development-related libraries are included (e.g., air for hot code reloading).
   - Uses Dockerfile and docker-compose.yml.
   - Follow these steps:
     - Start a PostgreSQL container:
       ```bash
       docker-compose up -d postgres
       ```
     - Build the project image:
       ```bash
       docker build -t go-server .
       ```
     - Run the built image:
       ```bash
       docker run -p 3000:3000 go-server
       ```